package main

import (
	"container/list"
	"fmt"
	"strconv"
	"strings"
)

type treeNode struct {
	parent   string // LHSsym
	termsym  string // must be +, *, (, ), id
	children *treeNodeQueue
}

type treeNodeQueue struct { // has one field, which is a ptr to a list
	queue *list.List // list elements don’t have datatype declaration
}

type treeStack struct {
	stack *list.List
}

var grammar = [6][]string{
	{"E", "->", "E", "+", "T"}, // 0
	{"E", "->", "T"},           // 1
	{"T", "->", "T", "*", "F"}, // 2
	{"T", "->", "F"},           // 3
	{"F", "->", "(", "E", ")"}, // 4
	{"F", "->", "id"}}          // 5

var aTable = [12][6]string{ // action table
	{"S5", "", "", "S4", "", ""},     // 0
	{"", "S6", "", "", "", "accept"}, // 1
	{"", "R2", "S7", "", "R2", "R2"}, // 2
	{"", "R4", "R4", "", "R4", "R4"}, // 3
	{"S5", "", "", "S4", "", ""},     // 4
	{"", "R6", "R6", "", "R6", "R6"}, // 5
	{"S5", "", "", "S4", "", ""},     // 6
	{"S5", "", "", "S4", "", ""},     // 7
	{"", "S6", "", "", "S11", ""},    // 8
	{"", "R1", "S7", "", "R1", "R1"}, // 9
	{"", "R3", "R3", "", "R3", "R3"}, // 10
	{"", "R5", "R5", "", "R5", "R5"}, // 11
}

var gTable = [12][3]string{
	{"1", "2", "3"}, // 0
	{"", "", ""},    // 1
	{"", "", ""},    // 2
	{"", "", ""},    // 3
	{"8", "2", "3"}, // 4
	{"", "", ""},    // 5
	{"", "9", "3"},  // 6
	{"", "", "10"},  // 7
	{"", "", ""},    // 8
	{"", "", ""},    // 9
	{"", "", ""},    // 10
	{"", "", ""},    // 11
}

var notGrammatical = false
var accept = "accept"
var ungrammatical = ""
var shift = "S"
var reduce = "R"

func enqueue(queue []string, element string) []string {
	queue = append(queue, element)
	return queue
}

func dequeue(queue []string) (string, []string) {
	element := queue[0]
	if len(queue) == 1 {
		var tmp = []string{}
		return element, tmp
	}
	return element, queue[1:]
}

type pstackItem struct {
	grammarSym, stateSym string
}

func (se pstackItem) String() string {
	return se.grammarSym + se.stateSym
}

type parseStack struct {
	stack *list.List
}

func (stk parseStack) String() string {
	stackString := ""
	for e := stk.stack.Front(); e != nil; e = e.Next() {
		stackString = (e.Value.(pstackItem)).String() + stackString
	}
	return stackString
}

// Creates an empty list to be a stack
func newParseStack() parseStack {
	ps := parseStack{}
	ps.stack = list.New()
	return ps
}

func (stk parseStack) push(itm pstackItem) {
	stk.stack.PushFront(itm)
}

func (stk parseStack) top() pstackItem {
	e := stk.stack.Front()
	return e.Value.(pstackItem)
}

func (stk parseStack) pop() pstackItem {
	e := stk.stack.Front()
	if e != nil {
		stk.stack.Remove(e)
		return e.Value.(pstackItem)
	}
	return pstackItem{"", ""}
}

func (stk parseStack) popNum(n int) {
	for i := 0; i < n; i++ {
		stk.pop()
	}
}

func handleInput(s string) int {
	if s == "id" {
		return 0
	} else if s == "+" {
		return 1
	} else if s == "*" {
		return 2
	} else if s == "(" {
		return 3
	} else if s == ")" {
		return 4
	} else if s == "$" {
		return 5
	} else {
		return -1 //MIGHT BE A PROBLEM HERE
	}
}

func handleGOTO(s string) int {
	if s == "E" {
		return 0
	} else if s == "T" {
		return 1
	} else if s == "F" {
		return 2
	} else {
		return -1 //MIGHT BE A PROBLEM HERE
	}
}

func determineR(i int) int { // MIGHT BE A PROBLEM IN THIS FUNCTION
	rule := grammar[i-1]
	ruleSlice := rule[2:]
	ruleLength := len(ruleSlice) - 1
	return ruleLength
}

func determineLHS(i int) string {
	rule := grammar[i-1]
	ruleLHS := rule[0]
	return ruleLHS
}

func parse1step(input []string) {
	/*switch choice {
	case accept:
		break
	case ungrammatical:
		notGrammatical = true
		break
	case shift:
	case reduce:
	}*/

	stack := newParseStack()

	stack.push(pstackItem{"", "0"})

	//current_token, input := dequeue(input)
	current_token := input[0]
	currentState := stack.top().stateSym
	cS, _ := strconv.Atoi(currentState)
	actionValue := aTable[cS][handleInput(current_token)]
	fmt.Printf("%-14s %-14s [%2d,%2s]  %-6s\n", stack.String(), input, cS, current_token, actionValue)
	for {
		if strings.Split(actionValue, "")[0] == shift {
			stack.push(pstackItem{current_token, strings.Split(actionValue, "")[1]})
			current_token, input = dequeue(input)
			fmt.Printf("%-14s %-14s [%2d,%2s]  %-6s\n", stack.String(), input, cS, current_token, actionValue)
		} else if strings.Split(actionValue, "")[0] == reduce {
			grammarNum := strings.Split(actionValue, "")[1]
			gN, _ := strconv.Atoi(grammarNum)
			stack.popNum(determineR(gN)) //NOT SURE ABOUT THIS
			currentState = stack.top().stateSym
			cS, _ = strconv.Atoi(currentState)
			stack.push(pstackItem{currentState, gTable[cS][handleGOTO(determineLHS(gN))]})
			fmt.Printf("%-14s %-14s [%2d,%2s]  %-6s\n      %s       %d               [%d,%s]     %s\n", stack.String(), input, cS, current_token, actionValue, determineLHS(gN), determineR(gN), cS, determineLHS(gN), gTable[cS][handleGOTO(determineLHS(gN))])
		} else if actionValue == ungrammatical { //MIGHT BE AN ERROR
			notGrammatical = true
			break
		} else if actionValue == accept {
			fmt.Printf("%-14s %-14s [%2d,%2s]  %-6s\n", stack.String(), input, cS, current_token, actionValue)
			break
		}
	}
}

func main() {
	inputArray := []string{"id", "+", "id", "*", "id"}
	//inputArray = []string{"(", "id", ")"}

	//stack := newParseStack()

	//stack.push(pstackItem{"0", ""})

	var inputQueue = make([]string, 0)

	for i := 0; i < len(inputArray); i++ {
		inputQueue = enqueue(inputQueue, inputArray[i])
	}
	inputQueue = enqueue(inputQueue, "$")

	fmt.Println("                input          action    action  value   length  temp            goto      goto   stack")
	fmt.Println("Stack           tokens         lookup    value   of LHS  of RHS  stack           lookup    value  action      parse tree stack")
	fmt.Println("______________________________________________________________________________________________________________________________")
	parse1step(inputQueue)
}

/* MIGHT BE USEFUL
push(0);
read_next_token();
for(;;)
{  s = top();    // current state is taken from top of stack
	if (ACTION[s,current_token] == 'si')   // shift and go to state i
	{  push(i);
	   read_next_token();
	}
	else if (ACTION[s,current_token] == 'ri')
	// reduce by rule i: X ::= A1...An
	{  perform pop() n times;
	   s = top();    // restore state before reduction from top of stack
	   push(GOTO[s,X]);   // state after reduction
	}
	else if (ACTION[s,current_token] == 'a')
	   success!!
	else error();
 }
*/
