package main

import (
    "bytes"
    "context"
    "fmt"
    "io"
    "log"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "sync"
    "time"

    "golang.org/x/text/encoding/json"
    "github.com/mattn/go-isatty"
    "github.com/rivo/uniseg"
    "github.com/charmbraceful/anker"
    "github.com/charmbraceful/anker/ankerui"
    "github.com/charmbraceful/anker/args"
    "github.com/charmbraceful/anker/terminal"
    "github.com/charmbraceful/anker/widgets"
)

type Config struct {
    Tools struct {
        Recon    map[string]string `json:"recon"`
        Scanning map[string]string `json:"scanning"`
        Crawling map[string]string `json:"crawling"`
        Exploit  map[string]string `json:"exploit"`
    } `json:"tools"`
    BBot struct {
        Target string `json:"target"`
    } `json:"bbot"`
}

type TargetConfig struct {
    Domains []string `json:"domains"`
}

// Task represents a task that can be executed
type Task struct {
    Name       string
    Command    string
    Args       string
    Category   string
    Output     bytes.Buffer
    Error      error
    OutputDir  string
    TargetDomains []string
    IsPiping   bool
    PipeTarget *Task
}

// TaskResult holds the result of a task execution
type TaskResult struct {
    Output   bytes.Buffer
    Error    error
    TaskName string
}

func main() {
    if !isatty.IsTerminal(os.Stdout.Fd()) {
        log.Println("not a terminal")
        return
    }

    // Initialize the application
    app := ankerui.New()
    app.SetTheme(ankerui.DefaultTheme())

    // Load configuration files
    cfg, err := loadConfig()
    if err != nil {
        log.Fatal(err)
    }

    tgtCfg, err := loadTargetConfig()
    if err != nil {
        log.Fatal(err)
    }

    // Initialize UI components
    listbox := createTaskList(cfg)
    console := createConsoleWidget()
    argsTextbox := createArgsTextbox()
    domainInput := createDomainInput()
    outputDirInput := createOutputDirInput()
    pipeMode := false

    // Create a channel for task results
    resultChan := make(chan TaskResult, 10)

    // Start button with task execution logic
    startButton := createStartButton(listbox, argsTextbox, domainInput, outputDirInput, resultChan, &pipeMode)
    pipeButton := createPipeButton(&pipeMode, listbox, argsTextbox, domainInput, outputDirInput, resultChan)

    // Create a layout
    layout := widgets.NewFlex().
        AddWidget(listbox, widgets.FlexSpec{Width: 25}).
        AddWidget(widgets.NewVerticalScroll(
            widgets.NewFlex().
                AddWidget(startButton, widgets.FlexSpec{Width: 20}).
                AddWidget(pipeButton, widgets.FlexSpec{Width: 20}).
                AddWidget(argsTextbox, widgets.FlexSpec{Width: 30}).
                AddWidget(domainInput, widgets.FlexSpec{Width: 20}).
                AddWidget(outputDirInput, widgets.FlexSpec{Width: 20}).
                AddWidget(console, widgets.FlexSpec{Width: 100}),
        ))

    // Set up the main view
    app.SetRoot(layout)

    // Start the application
    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}

func loadConfig() (*Config, error) {
    configFile, err := filepath.Abs("./config.json")
    if err != nil {
        return nil, fmt.Errorf("failed to get the absolute path of config.json: %v", err)
    }

    file, err := os.Open(configFile)
    if err != nil {
        return nil, fmt.Errorf("failed to open config.json: %v", err)
    }
    defer file.Close()

    var cfg Config
    if err := json.NewDecoder(file).Decode(&cfg); err != nil {
        return nil, fmt.Errorf("failed to decode the configuration file: %v", err)
    }

    return &cfg, nil
}

func loadTargetConfig() (*TargetConfig, error) {
    targetConfigFile, err := filepath.Abs("./target_domains.json")
    if err != nil {
        return nil, fmt.Errorf("failed to get the absolute path of target_domains.json: %v", err)
    }

    file, err := os.Open(targetConfigFile)
    if err != nil {
        return nil, fmt.Errorf("failed to open target_domains.json: %v", err)
    }
    defer file.Close()

    var tgtCfg TargetConfig
    if err := json.NewDecoder(file).Decode(&tgtCfg); err != nil {
        return nil, fmt.Errorf("failed to decode the target_domains.json file: %v", err)
    }

    return &tgtCfg, nil
}

func createTaskList(cfg *Config) *widgets.ListBox {
    items := make([]string, 0)

    // Add tasks with categories
    for category, tools := range cfg.Tools {
        for name, _ := range tools {
            items = append(items, fmt.Sprintf("[%s] %s", category, name))
        }
    }

    listbox, err := widgets.NewListBox(
        widgets.NewAnkerUI(),
        widgets.ListboxOptions{
            Items: items,
        },
    )
    if err != nil {
        log.Fatal(err)
    }
    listbox.SetTheme(ankerui.DefaultTheme())

    return listbox
}

func createConsoleWidget() *widgets.Textarea {
    console, err := widgets.NewTextarea(
        widgets.NewAnkerUI(),
        widgets.TextareaOptions{
            Rows: 15,
        },
    )
    if err != nil {
        log.Fatal(err)
    }
    console.SetTheme(ankerui.DefaultTheme())
    return console
}

func createArgsTextbox() *widgets.Textbox {
    argsTextbox, err := widgets.NewTextbox(
        widgets.NewAnkerUI(),
        widgets.TextboxOptions{
            Rows: 1,
            Title: "Custom Arguments",
        },
    )
    if err != nil {
        log.Fatal(err)
    }
    argsTextbox.SetTheme(ankerui.DefaultTheme())
    return argsTextbox
}

func createDomainInput() *widgets.Textbox {
    domainInput, err := widgets.NewTextbox(
        widgets.NewAnkerUI(),
        widgets.TextboxOptions{
            Rows: 1,
            Title: "Target Domain(s)",
        },
    )
    if err != nil {
        log.Fatal(err)
    }
    domainInput.SetTheme(ankerui.DefaultTheme())
    return domainInput
}

func createOutputDirInput() *widgets.Textbox {
    outputDirInput, err := widgets.NewTextbox(
        widgets.NewAnkerUI(),
        widgets.TextboxOptions{
            Rows: 1,
            Title: "Output Directory",
        },
    )
    if err != nil {
        log.Fatal(err)
    }
    outputDirInput.SetTheme(ankerui.DefaultTheme())
    return outputDirInput
}

func createStartButton(listbox *widgets.ListBox, argsTextbox *widgets.Textbox, domainInput *widgets.Textbox, outputDirInput *widgets.Textbox, resultChan chan TaskResult, pipeMode *bool) *widgets.Button {
    startButton := widgets.NewButton(
        "Start",
        func(e anker.Event) {
            // Get selected task
            taskName := listbox.SelectedItem()
            if taskName == "" {
                return
            }

            // Extract task name without category
            task := strings.Split(taskName, "] ")[1]

            // Get custom arguments
            args := argsTextbox.Value()

            // Get target domain(s)
            domains := domainInput.Value()

            // Get output directory
            outputDir := outputDirInput.Value()

            // Validate input
            if err := validateInput(domains, outputDir); err != nil {
                console.SetValue(fmt.Sprintf("Error: %v", err))
                return
            }

            // Prepare the task
            t := prepareTask(task, domains, args, outputDir, *pipeMode, nil)

            // Execute the task asynchronously
            go executeTask(t, resultChan)
        },
    )
    startButton.SetTheme(ankerui.DefaultTheme())
    return startButton
}

func createPipeButton(pipeMode *bool, listbox *widgets.ListBox, argsTextbox *widgets.Textbox, domainInput *widgets.Textbox, outputDirInput *widgets.Textbox, resultChan chan TaskResult) *widgets.Button {
    pipeButton := widgets.NewButton(
        "Pipe",
        func(e anker.Event) {
            *pipeMode = !*pipeMode
            if *pipeMode {
                // Prepare for piping
                // Select first tool
                taskName := listbox.SelectedItem()
                if taskName == "" {
                    *pipeMode = false
                    return
                }

                // Extract task name without category
                task := strings.Split(taskName, "] ")[1]

                // Get custom arguments
                args := argsTextbox.Value()

                // Get target domain(s)
                domains := domainInput.Value()

                // Get output directory
                outputDir := outputDirInput.Value()

                // Validate input
                if err := validateInput(domains, outputDir); err != nil {
                    console.SetValue(fmt.Sprintf("Error: %v", err))
                    *pipeMode = false
                    return
                }

                // Prepare the first task
                t := prepareTask(task, domains, args, outputDir, *pipeMode, nil)

                // Execute the first task and prepare for the next
                go executeTask(t, resultChan)
            }
        },
    )
    pipeButton.SetTheme(ankerui.DefaultTheme())
    return pipeButton
}

func validateInput(domains, outputDir string) error {
    if domains == "" {
        return fmt.Errorf("target domain cannot be empty")
    }
    if outputDir == "" {
        return fmt.Errorf("output directory cannot be empty")
    }
    if _, err := os.Stat(outputDir); os.IsNotExist(err) {
        return fmt.Errorf("output directory does not exist")
    }
    return nil
}

func prepareTask(taskName, domains, args, outputDir string, isPiping bool, pipeTarget *Task) *Task {
    // Split domains by comma
    domainList := strings.Split(domains, ",")
    
    // Clean up domains
    for i := range domainList {
        domainList[i] = strings.TrimSpace(domainList[i])
    }

    // Prepare arguments for each domain
    argsList := make([]string, 0)
    for _, domain := range domainList {
        argsList = append(argsList, fmt.Sprintf("--domain %s", domain))
    }

    // Combine all arguments
    allArgs := strings.Join(argsList, " ") + " " + args

    return &Task{
        Name:       taskName,
        Command:    getCommandPath(taskName),
        Args:       allArgs,
        Category:   strings.Split(taskName, "] ")[0][1:],
        Output:     *bytes.NewBuffer([]byte{}),
        Error:      nil,
        OutputDir:  outputDir,
        TargetDomains: domainList,
        IsPiping:   isPiping,
        PipeTarget: pipeTarget,
    }
}

func getCommandPath(taskName string) string {
    // Implementation to get the command path from config
    // This is a simplified version
    return taskName
}

func executeTask(t *Task, resultChan chan TaskResult) {
    ctx := context.Background()

    cmd := exec.CommandContext(ctx, t.Command, strings.Split(t.Args, " ")...)
    cmd.Stdout = &t.Output
    cmd.Stderr = &t.Output

    var wg sync.WaitGroup
    wg.Add(1)

    go func() {
        defer wg.Done()
        if err := cmd.Run(); err != nil {
            t.Error = fmt.Errorf("command failed: %v", err)
        }
    }()

    // Stream output in real-time
    go func() {
        buf := make([]byte, 1024)
        for {
            n, err := cmd.Stdout.Read(buf)
            if err == io.EOF {
                break
            }
            if err != nil {
                t.Error = fmt.Errorf("error reading output: %v", err)
                break
            }
            t.Output.Write(buf[:n])
            resultChan <- *t
        }
    }()

    wg.Wait()

    resultChan <- *t
}
