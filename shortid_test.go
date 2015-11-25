// Copyright 2015 Bryan Weber. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// The shortid_test package provides tests for the shortid package
package shortid_test

import (
	"github.com/bradialabs/shortid"
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {
	s := shortid.New()
	s.SetSeed(12345)
	s.SetWorkerId(6)
	i := 5000
	ids := make(map[string]int)

	for i > 0 {
		id := s.Generate()
		t.Logf("Shortened id: %s", id)
		_, exists := ids[id]
		if exists {
			ids[id] = ids[id] + 1
		} else {
			ids[id] = 1
		}

		i = i - 1

		if ids[id] > 1 {
			t.Errorf("Generated non-unique IDs: %s - %d (%d)", id, ids[id], time.Now().Unix())
		}
	}
}

func TestDecode(t *testing.T) {
	s := shortid.New()
	s.SetWorkerId(0)
	_, worker := s.Decode(s.Generate())
	if worker != 0 {
		t.Errorf("Worker not decoded correctly")
	}

	s.SetWorkerId(1)
	_, worker = s.Decode(s.Generate())
	if worker != 1 {
		t.Errorf("Worker not decoded correctly")
	}

	s.SetWorkerId(15)
	_, worker = s.Decode(s.Generate())
	if worker != 15 {
		t.Errorf("Worker not decoded correctly")
	}
}
