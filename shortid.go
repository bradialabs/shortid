package shortid

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"math"
	"strings"
	"time"
)

//reduce the unix seconds by this amount to keep IDs small.
//Update about yearly and increment the VERSION to avoid collisions
const REDUCE_TIME int64 = 1448403506

//Only change the version when REDUCE_TIME is updated or an algorithm change.
//Must be an integer below 16.
const VERSION int64 = 1

// characters used for conversion
const ORIGINAL = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_-"

type ShortId struct {
	clusterWorkerId int64
	counter         int64
	prevSeconds     int64
	seed            int64
	alphabet        string
	shuffled        string
}

func New() *ShortId {
	s := &ShortId{
		clusterWorkerId: 0,        //worker or machine id (positive int less than 16)
		counter:         0,        //internal counter for generation in the same second
		prevSeconds:     0,        //unix seconds of previous call
		seed:            1,        //number to seed the shuffling of the alphabet
		alphabet:        ORIGINAL, //the unshuffled alphabet
		shuffled:        "",       //the shuffled alphabet used for encoding
	}

	return s
}

/******************************************
* Generates a new unique short id.
 */
func (s *ShortId) Generate() string {
	str := ""
	//get the current unix seconds
	seconds := time.Now().Unix() - REDUCE_TIME

	if seconds == s.prevSeconds {
		s.counter = s.counter + 1
	} else {
		s.counter = 0
		s.prevSeconds = seconds
	}

	str = str + s.encode(VERSION)
	str = str + s.encode(s.clusterWorkerId)
	if s.counter > 0 {
		str = str + s.encode(s.counter)
	}
	str = str + s.encode(seconds)

	return str
}

/*****************************************************
* Encodes an int into a character from the shuffled alphabet
 */
func (s *ShortId) encode(number int64) string {
	loopCounter := int64(0)
	done := false
	b := make([]byte, 1)
	rand.Read(b)
	buf := bytes.NewBuffer(b)
	r, _ := binary.ReadVarint(buf)
	r = r & 0x30

	str := ""

	for done != true {
		lookupIndex := ((number >> uint64(4*loopCounter)) & 0x0f) | r
		str = str + s.lookup(lookupIndex)
		done = float64(number) < (math.Pow(16, float64(loopCounter+1)))
		loopCounter = loopCounter + 1
	}
	return str
}

/**********************************************
* Decodes a short id to get the worker and version.
* Used mostly just for debugging.
 */
func (s *ShortId) Decode(id string) (version, worker int) {
	characters := s.getShuffled()
	version = strings.Index(characters, id[0:1]) & 0x0f
	worker = strings.Index(characters, id[1:2]) & 0x0f
	return
}

/***********************************************
* Found this seed-based random generator somewhere
* Based on The Central Randomizer 1.3 (C) 1997 by Paul Houle (houle@msc.cornell.edu)
 */
func (s *ShortId) getRandomValue() float64 {
	s.seed = (s.seed*9301 + 49297) % 233280
	return float64(s.seed) / 233280.0
}

/************************************************
* Shuffle the alphabet using the seeded random
* number generator.
 */
func (s *ShortId) shuffle() string {
	source := []byte(s.alphabet)
	var target []byte
	r := s.getRandomValue()
	var charIndex int64

	for len(source) > 0 {
		r = s.getRandomValue()
		charIndex = int64(math.Floor(r * float64(len(source))))
		target = append(target, source[charIndex])
		if charIndex < int64(len(source)) {
			source = append(source[:charIndex], source[charIndex+1:]...)
		} else {
			source = source[:charIndex]
		}
	}

	return string(target)
}

/************************************************
* Get the shuffled alphabet, and shuffle it if it
* hasn't been already.
 */
func (s *ShortId) getShuffled() string {
	if len(s.shuffled) > 0 {
		return s.shuffled
	}
	s.shuffled = s.shuffle()
	return s.shuffled
}

/************************************************
* Get a character from the shuffled alphabet
 */
func (s *ShortId) lookup(index int64) string {
	alphabetShuffled := s.getShuffled()
	return string(alphabetShuffled[index])
}

/************************************************
* Set the seed for shuffling
 */
func (s *ShortId) SetSeed(seed int64) {
	s.seed = seed
}

/************************************************
* Set the worker id
 */
func (s *ShortId) SetWorkerId(workerId int64) {
	s.clusterWorkerId = workerId
}

/************************************************
* Set the alphabet
 */
func (s *ShortId) SetAlphabet(alpha string) {
	s.alphabet = alpha
}
