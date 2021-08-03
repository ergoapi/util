//  Copyright (c) 2020. The EFF Team Authors.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  See the License for the specific language governing permissions and
//  limitations under the License.

package ptr

// IntPtr int指针
func IntPtr(p int) *int { return &p }

// Int32Ptr int32指针
func Int32Ptr(p int32) *int32 { return &p }

// Int64Ptr int64指针
func Int64Ptr(p int64) *int64 { return &p }

// StringPtr 字符串指针
func StringPtr(p string) *string { return &p }

// BoolPtr bool指针
func BoolPtr(p bool) *bool { return &p }
