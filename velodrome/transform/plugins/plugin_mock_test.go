/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Automatically generated by MockGen. DO NOT EDIT!
// Source: plugin.go

package plugins

import (
	gomock "github.com/golang/mock/gomock"
	sql "k8s.io/test-infra/velodrome/sql"
)

// Mock of Plugin interface
type MockPlugin struct {
	ctrl     *gomock.Controller
	recorder *MockPluginRecorder
}

// Recorder for MockPlugin
type MockPluginRecorder struct {
	mock *MockPlugin
}

func NewMockPlugin(ctrl *gomock.Controller) *MockPlugin {
	mock := &MockPlugin{ctrl: ctrl}
	mock.recorder = &MockPluginRecorder{mock}
	return mock
}

func (_m *MockPlugin) EXPECT() *MockPluginRecorder {
	return _m.recorder
}

func (_m *MockPlugin) ReceiveIssue(_param0 sql.Issue) []Point {
	ret := _m.ctrl.Call(_m, "ReceiveIssue", _param0)
	ret0, _ := ret[0].([]Point)
	return ret0
}

func (_mr *MockPluginRecorder) ReceiveIssue(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ReceiveIssue", arg0)
}

func (_m *MockPlugin) ReceiveComment(_param0 sql.Comment) []Point {
	ret := _m.ctrl.Call(_m, "ReceiveComment", _param0)
	ret0, _ := ret[0].([]Point)
	return ret0
}

func (_mr *MockPluginRecorder) ReceiveComment(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ReceiveComment", arg0)
}

func (_m *MockPlugin) ReceiveIssueEvent(_param0 sql.IssueEvent) []Point {
	ret := _m.ctrl.Call(_m, "ReceiveIssueEvent", _param0)
	ret0, _ := ret[0].([]Point)
	return ret0
}

func (_mr *MockPluginRecorder) ReceiveIssueEvent(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ReceiveIssueEvent", arg0)
}