package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// ─── Colors ───────────────────────────────────────────────────────────────────

const (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Dim    = "\033[2m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
	White  = "\033[97m"
)

// ─── Data Structures ──────────────────────────────────────────────────────────

type Priority string

const (
	Low    Priority = "low"
	Medium Priority = "medium"
	High   Priority = "high"
)

type Task struct {
	ID        int      `json:"id"`
	Title     string   `json:"title"`
	Done      bool     `json:"done"`
	Priority  Priority `json:"priority"`
	CreatedAt string   `json:"created_at"`
	DoneAt    string   `json:"done_at,omitempty"`
}

type Store struct {
	Tasks   []Task `json:"tasks"`
	Counter int    `json:"counter"`
}

const dataFile = "tasks.json"

// ─── Storage ──────────────────────────────────────────────────────────────────

func loadStore() Store {
	data, err := os.ReadFile(dataFile)
	if err != nil {
		return Store{}
	}
	var s Store
	json.Unmarshal(data, &s)
	return s
}

func saveStore(s Store) {
	data, _ := json.MarshalIndent(s, "", "  ")
	os.WriteFile(dataFile, data, 0644)
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func priorityColor(p Priority) string {
	switch p {
	case High:
		return Red
	case Medium:
		return Yellow
	default:
		return Dim
	}
}

func priorityIcon(p Priority) string {
	switch p {
	case High:
		return "!!!"
	case Medium:
		return "! "
	default:
		return "  "
	}
}

func parsePriority(s string) Priority {
	switch strings.ToLower(s) {
	case "h", "high":
		return High
	case "m", "medium":
		return Medium
	default:
		return Low
	}
}

func now() string {
	return time.Now().Format("2006-01-02 15:04")
}

func printHeader() {
	fmt.Println()
	fmt.Printf("%s%s╔══════════════════════════════╗%s\n", Bold, Cyan, Reset)
	fmt.Printf("%s%s║     ✅  Go Task Manager      ║%s\n", Bold, Cyan, Reset)
	fmt.Printf("%s%s╚══════════════════════════════╝%s\n", Bold, Cyan, Reset)
	fmt.Println()
}

func printHelp() {
	printHeader()
	fmt.Printf("%sUSAGE:%s\n", Bold, Reset)
	fmt.Printf("  %stask%s <command> [arguments]\n\n", Cyan, Reset)

	fmt.Printf("%sCOMMANDS:%s\n", Bold, Reset)
	cmds := [][]string{
		{"add <title> [priority]", "Add a new task  (priority: low/medium/high, default: low)"},
		{"list", "List all tasks"},
		{"list pending", "List only pending tasks"},
		{"list done", "List only completed tasks"},
		{"done <id>", "Mark a task as complete"},
		{"undone <id>", "Mark a task as pending again"},
		{"delete <id>", "Delete a task"},
		{"clear done", "Remove all completed tasks"},
		{"stats", "Show summary statistics"},
		{"help", "Show this help message"},
	}

	for _, c := range cmds {
		fmt.Printf("  %s%-30s%s %s%s%s\n", Green, c[0], Reset, Dim, c[1], Reset)
	}

	fmt.Printf("\n%sEXAMPLES:%s\n", Bold, Reset)
	fmt.Printf("  %sgo run main.go add \"Write unit tests\" high%s\n", Dim, Reset)
	fmt.Printf("  %sgo run main.go done 3%s\n", Dim, Reset)
	fmt.Printf("  %sgo run main.go list pending%s\n\n", Dim, Reset)
}

// ─── Commands ─────────────────────────────────────────────────────────────────

func cmdAdd(args []string) {
	if len(args) == 0 {
		fmt.Printf("%s✗ Usage: add <title> [priority]%s\n", Red, Reset)
		return
	}

	title := args[0]
	priority := Low
	if len(args) > 1 {
		priority = parsePriority(args[1])
	}

	s := loadStore()
	s.Counter++
	task := Task{
		ID:        s.Counter,
		Title:     title,
		Done:      false,
		Priority:  priority,
		CreatedAt: now(),
	}
	s.Tasks = append(s.Tasks, task)
	saveStore(s)

	pc := priorityColor(priority)
	fmt.Printf("%s✓ Task added:%s [#%d] %s %s(%s)%s\n",
		Green, Reset, task.ID, task.Title, pc, string(priority), Reset)
}

func cmdList(filter string) {
	s := loadStore()
	if len(s.Tasks) == 0 {
		fmt.Printf("%s  No tasks yet. Add one with: add <title>%s\n\n", Dim, Reset)
		return
	}

	var pending, done []Task
	for _, t := range s.Tasks {
		if t.Done {
			done = append(done, t)
		} else {
			pending = append(pending, t)
		}
	}

	printHeader()

	showPending := filter == "" || filter == "pending"
	showDone := filter == "" || filter == "done"

	if showPending {
		fmt.Printf("%s%s  PENDING  (%d)%s\n", Bold, Yellow, len(pending), Reset)
		fmt.Printf("%s  ─────────────────────────────────────────%s\n", Dim, Reset)
		if len(pending) == 0 {
			fmt.Printf("%s  All done! Nothing pending 🎉%s\n", Dim, Reset)
		}
		for _, t := range pending {
			pc := priorityColor(t.Priority)
			icon := priorityIcon(t.Priority)
			fmt.Printf("  %s[%s]%s  %s%s%s  %s  %s%s%s\n",
				Dim, fmt.Sprintf("%2d", t.ID), Reset,
				pc, icon, Reset,
				t.Title,
				Dim, t.CreatedAt, Reset)
		}
		fmt.Println()
	}

	if showDone {
		fmt.Printf("%s%s  COMPLETED  (%d)%s\n", Bold, Green, len(done), Reset)
		fmt.Printf("%s  ─────────────────────────────────────────%s\n", Dim, Reset)
		if len(done) == 0 {
			fmt.Printf("%s  No completed tasks yet.%s\n", Dim, Reset)
		}
		for _, t := range done {
			fmt.Printf("  %s[%2d]%s  %s✓%s  %s%s%s  %s%s%s\n",
				Dim, t.ID, Reset,
				Green, Reset,
				Dim, t.Title, Reset,
				Dim, t.DoneAt, Reset)
		}
		fmt.Println()
	}
}

func findTask(s *Store, id int) int {
	for i, t := range s.Tasks {
		if t.ID == id {
			return i
		}
	}
	return -1
}

func cmdDone(idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Printf("%s✗ Invalid ID: %s%s\n", Red, idStr, Reset)
		return
	}

	s := loadStore()
	idx := findTask(&s, id)
	if idx == -1 {
		fmt.Printf("%s✗ Task #%d not found.%s\n", Red, id, Reset)
		return
	}
	if s.Tasks[idx].Done {
		fmt.Printf("%s  Task #%d is already done.%s\n", Dim, id, Reset)
		return
	}

	s.Tasks[idx].Done = true
	s.Tasks[idx].DoneAt = now()
	saveStore(s)
	fmt.Printf("%s✓ Task #%d marked as done:%s %s\n", Green, id, Reset, s.Tasks[idx].Title)
}

func cmdUndone(idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Printf("%s✗ Invalid ID: %s%s\n", Red, idStr, Reset)
		return
	}

	s := loadStore()
	idx := findTask(&s, id)
	if idx == -1 {
		fmt.Printf("%s✗ Task #%d not found.%s\n", Red, id, Reset)
		return
	}

	s.Tasks[idx].Done = false
	s.Tasks[idx].DoneAt = ""
	saveStore(s)
	fmt.Printf("%s↩ Task #%d moved back to pending:%s %s\n", Yellow, id, Reset, s.Tasks[idx].Title)
}

func cmdDelete(idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Printf("%s✗ Invalid ID: %s%s\n", Red, idStr, Reset)
		return
	}

	s := loadStore()
	idx := findTask(&s, id)
	if idx == -1 {
		fmt.Printf("%s✗ Task #%d not found.%s\n", Red, id, Reset)
		return
	}

	title := s.Tasks[idx].Title
	s.Tasks = append(s.Tasks[:idx], s.Tasks[idx+1:]...)
	saveStore(s)
	fmt.Printf("%s✗ Deleted task #%d:%s %s\n", Red, id, Reset, title)
}

func cmdClearDone() {
	s := loadStore()
	var remaining []Task
	count := 0
	for _, t := range s.Tasks {
		if !t.Done {
			remaining = append(remaining, t)
		} else {
			count++
		}
	}
	s.Tasks = remaining
	saveStore(s)
	fmt.Printf("%s✓ Cleared %d completed task(s).%s\n", Green, count, Reset)
}

func cmdStats() {
	s := loadStore()
	total := len(s.Tasks)
	done := 0
	byPriority := map[Priority]int{Low: 0, Medium: 0, High: 0}

	for _, t := range s.Tasks {
		if t.Done {
			done++
		}
		if !t.Done {
			byPriority[t.Priority]++
		}
	}
	pending := total - done

	var pct int
	if total > 0 {
		pct = (done * 100) / total
	}

	// Progress bar
	bar := 20
	filled := (pct * bar) / 100
	progress := strings.Repeat("█", filled) + strings.Repeat("░", bar-filled)

	printHeader()
	fmt.Printf("%s  STATISTICS%s\n\n", Bold, Reset)
	fmt.Printf("  Total Tasks   : %s%d%s\n", Bold, total, Reset)
	fmt.Printf("  Pending       : %s%d%s\n", Yellow, pending, Reset)
	fmt.Printf("  Completed     : %s%d%s\n", Green, done, Reset)
	fmt.Printf("  Progress      : %s%s%s %s%d%%%s\n\n", Cyan, progress, Reset, Bold, pct, Reset)
	fmt.Printf("  Pending by Priority:\n")
	fmt.Printf("    %s!!! High   : %d%s\n", Red, byPriority[High], Reset)
	fmt.Printf("    %s!   Medium : %d%s\n", Yellow, byPriority[Medium], Reset)
	fmt.Printf("    %s    Low    : %d%s\n\n", Dim, byPriority[Low], Reset)
}

// ─── Main ─────────────────────────────────────────────────────────────────────

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		cmdList("")
		return
	}

	cmd := strings.ToLower(args[0])
	rest := args[1:]

	switch cmd {
	case "add":
		cmdAdd(rest)
	case "list":
		filter := ""
		if len(rest) > 0 {
			filter = strings.ToLower(rest[0])
		}
		cmdList(filter)
	case "done":
		if len(rest) == 0 {
			fmt.Printf("%s✗ Usage: done <id>%s\n", Red, Reset)
			return
		}
		cmdDone(rest[0])
	case "undone":
		if len(rest) == 0 {
			fmt.Printf("%s✗ Usage: undone <id>%s\n", Red, Reset)
			return
		}
		cmdUndone(rest[0])
	case "delete", "del", "rm":
		if len(rest) == 0 {
			fmt.Printf("%s✗ Usage: delete <id>%s\n", Red, Reset)
			return
		}
		cmdDelete(rest[0])
	case "clear":
		if len(rest) > 0 && rest[0] == "done" {
			cmdClearDone()
		} else {
			fmt.Printf("%s✗ Usage: clear done%s\n", Red, Reset)
		}
	case "stats":
		cmdStats()
	case "help", "--help", "-h":
		printHelp()
	default:
		fmt.Printf("%s✗ Unknown command: %s%s\n", Red, cmd, Reset)
		fmt.Printf("  Run %sgo run main.go help%s for usage.\n", Cyan, Reset)
	}
}
