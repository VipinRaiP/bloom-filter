package main

import (
  "fmt"
  "github.com/spaolacci/murmur3"
  "hash"
  "time"
  "github.com/google/uuid"  
)

type BloomFilter struct{
  filter []int8
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
    filter : make([]int8,size),
    size : size,
  }
}

func (b *BloomFilter) Add(key string){
    idx := murmurhash(key,b.size)
    aidx := idx/8
    bidx := idx % 8
    b.filter[aidx] = b.filter[aidx] | (1<<bidx)
}

func (b *BloomFilter) Exists(key string) bool{
  idx := murmurhash(key, b.size)
  aidx := idx/8
  bidx := idx%8
  return (b.filter[aidx] & (1<<bidx))>0
}

func main(){
  dataset := make([]string,1000)
  dataset_exists := make(map[string]bool)
  dataset_notexists := make(map[string]bool)

  for i:=1; i<=900; i++{
    key := uuid.New().String()
    dataset = append(dataset, key)
    dataset_exists[key] = true
  }

  for i:=1; i<=200; i++{
    key := uuid.New().String()
    dataset = append(dataset, key)
    dataset_notexists[key] = true
  }


  for i:=100;i<=10000;i+=100{
    bloom := NewBloomFilter(int32(i))
    for key,_ := range dataset_exists{
        bloom.Add(key)
    }
  
    false_positive_count := 0
    
    for _,key := range dataset{
      if bloom.Exists(key) && dataset_notexists[key]{
        false_positive_count++
      } 
    }  
  
    error_rate := 100 * (float64(false_positive_count)/float64(len(dataset_notexists)))
    fmt.Println(error_rate)
  }  
  
}
