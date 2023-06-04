package main

import (
  "fmt"
  "github.com/spaolacci/murmur3"
  "hash"
  "time"
  "github.com/google/uuid"  
)

type BloomFilter struct{
  filter []bool
  size int32
}

var mHasher hash.Hash32

func init(){
  mHasher = murmur3.New32WithSeed(uint32(time.Now().Unix()))  
}

func murmurhash(key string, size int32) uint32{
  mHasher.Write([]byte(key))
  result := mHasher.Sum32() % uint32(size)
  mHasher.Reset()
  return result
}

func NewBloomFilter(size int32) *BloomFilter{
  return &BloomFilter{
    filter : make([]bool,size),
    size : size,
  }
}

func (b *BloomFilter) Add(key string){
    idx := murmurhash(key,b.size)
    b.filter[idx] = true
}

func (b *BloomFilter) Exists(key string) bool{
  idx := murmurhash(key, b.size)
  return b.filter[idx]
}

func main(){
  dataset := make([]string,1000)
  dataset_exists := make(map[string]bool)
  dataset_notexists := make(map[string]bool)
  

  
  bloom := NewBloomFilter(16)
  keys := []string {"a","b","c"}
  for _,key := range keys{
      bloom.Add(key)
  }

  for _,key := range keys{
    fmt.Println(key, bloom.Exists(key))
  }  
}
