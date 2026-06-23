package navigator

import (
	"fmt"
	"strings"
)

type Level int
type ID string

//Levels: Root, System, Device, Sensor

const (
	Root = iota
	Zyztem
	Device
	Sensor
)

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
	rootNode := &NavigatorNode{ID: "0", Level: Root, Parent: nil, Children: make([]*NavigatorNode, 0)}
	navigator := &Navigator{CurrentNode: rootNode}
	return navigator, rootNode

}

func (parent *NavigatorNode) Add() (*NavigatorNode, error) {
	//TODO: make an iota for Level to easily do this
	if parent.Level == Sensor {
		return nil, fmt.Errorf("canot add child node to sensor")
	}
	childLevel := parent.Level + 1
	NewNode := &NavigatorNode{ID: "0", Level: childLevel, Parent: nil, Children: make([]*NavigatorNode, 0)}
	parent.Children = append(parent.Children, NewNode)
	return NewNode, nil
}

func (parent *NavigatorNode) Remove(child *NavigatorNode) error {
	//remove child
	//remove child from parent list
	//find match index in parent.Children and pop

	for i, c := range parent.Children {
		if c == child {
			parent.Children = append(parent.Children[:i], parent.Children[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("child node %q not found under parent %q", child.ID, parent.ID)
}

func (parent *NavigatorNode) List() string {
	if len(parent.Children) == 0 {
		return "(no children)"
	}
	var sb strings.Builder
	for i, c := range parent.Children {
		fmt.Fprintf(&sb, "%d: %s\n", i, c.ID)
	}
	return sb.String()

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
