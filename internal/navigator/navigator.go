package navigator

import "fmt"

type Level string
type ID string

//Levels: Root, System, Device, Sensor

type NavigatorNode struct {
	ID       ID //need to determine ID
	Level    Level
	Parent   *NavigatorNode
	Children []*NavigatorNode
}

type Navigator struct {
	CurrentNode *NavigatorNode
}

func New() (*Navigator, *NavigatorNode) {
	rootNode := &NavigatorNode{ID: "0", Level: "root", Parent: nil, Children: make([]*NavigatorNode, 0)}
	navigator := &Navigator{CurrentNode: rootNode}
	return navigator, rootNode

}

func (parent *NavigatorNode) Add() *NavigatorNode {
	//TODO: make an iota for Level to easily do this
	NewNode := &NavigatorNode{ID: "0", Level: "root", Parent: nil, Children: make([]*NavigatorNode, 0)}
	parent.Children = append(parent.Children, NewNode)
	return NewNode
}

func (parent *NavigatorNode) Remove(child *NavigatorNode) error {
	return fmt.Errorf("failed to remove child node")
}

func (parent *NavigatorNode) List() string {
	return "failed to Listchild node"
}

func (nav *Navigator) Set(chosenNode *NavigatorNode) {
	nav.CurrentNode = chosenNode
}

func (nav *Navigator) Up() error {
	return fmt.Errorf("failed to move navigator")
}

func (nav *Navigator) Down(id ID) error {
	return fmt.Errorf("failed to move navigator")
}

//Create functionality to read folder where systems should be (decide file format)
//Basically I think I'm just going to iterate the file names and store the PID then access all other information via this navigator (really good information to log to verify for clean-up)

//I imagine this will simply return the current node and directions

//nvigator will have a parent child tree system
