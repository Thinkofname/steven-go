// Copyright 2015 Matthew Collins
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package console implements a system for parsing and executing
// commands + logging.
//
// When registering commands a description of the command is
// required. For basic commands the format of this is simple:
//
//     commandname sub1 sub2 sub...
//
// Complex commands can be created by using % to specify
// arguments. The type of the argument will be inferred
// from the type over the passed function pointer. Extra
// constraints can be added after the % to gain finer control
// over the argument.
//
// Built-in types:
//     string        any string, a length limit may be added
//                   after the % to enforce a max length
//
// Executing commands works by treating whitespace at delimiters
// between arguments with the exception of whitespace contained
// within quotes (") as that will be treated as a single argument
package console
