package main

import (
    "fmt"
    "time"
    "runtime"
)

func fib(pos uint64) uint64 {
    if pos <= 1 {
        return pos;
    }
    return fib(pos - 1) + fib(pos - 2); 
}

func fibCached(pos uint64) uint64 {
    return fibWithCache(pos, make([]uint64, pos));
}

func fibWithCache(pos uint64, cache []uint64) uint64 {
    if pos <= 1 {
        return pos;
    }
    
    if cache[pos - 2] != 0 {
        return cache[pos - 2];
    }

    value := fibWithCache(pos - 1, cache) + fibWithCache(pos - 2, cache);
    if uint64(len(cache)) < pos {
        return value;
    }

    cache[pos - 2] = value;
    return value;
}

func fibCount(pos uint64) ([]uint64, []uint64) {
    count := make([]uint64, pos + 1);
    cache := make([]uint64, pos + 1);
    
    fibWithCacheAndCount(pos, cache, count); 
    return count, cache;
}
    
func fibWithCacheAndCount(pos uint64, cache []uint64, count []uint64) uint64 {
    if pos <= 1 {
        count[pos]++;
        cache[pos] = pos;
        return pos;
    }

    value := fibWithCacheAndCount(pos - 1, cache, count) + fibWithCacheAndCount(pos - 2, cache, count);
    count[pos]++;
    cache[pos] = value;
    return value;
}

func makeTest() {
    titles := []string{ "Fib Cached", "Fib jump", "Fib multi" };
    fn := []func(uint64)uint64{ fibCached, func(pos uint64) uint64 { a, _ := binFibJump(pos); return a } , multiFibJump };

    for i := range titles {
        fmt.Println("---", titles[i]);
        cache := make([]uint64, 201);

        start := time.Now();
        for j := range 201 {
            cache[j] = fn[i](uint64(j));
        }
    
        for i, c := range cache {
            fmt.Printf("fib(%d) = %d\n", i, c);
        } 
        fmt.Println("Tempo:", time.Now().Sub(start));
    }
}

func fibJump(pos uint64) uint64 {
    if pos <= 1 {
        return pos
    }

    const psize = 20;
    
    i, i0, i1 := fibIndexer(psize);
    f, f0 := i, i0;

    for iter := uint64(1); iter < pos / psize; iter++ {
        unit := f + f0;

        f0 = i0 * unit + i1 * f;
        f  = i  * unit + i0 * f;
    }

    if remain := pos % psize; remain != 0 {
        i, i0, i1 := fibIndexer(remain); 
        unit := f + f0;

        f0 = i0 * unit + i1 * f;
        f  = i  * unit + i0 * f;
    }
    return f;
}

func fibIndexer(pos uint64) (k, k0, k1 uint64) {
    k, k0, k1 = 1, 1, 0;

    if pos == 2 {
        return 1, 1, 0;
    }

    if pos == 1 {
        return 1, 0, 0;
    }

    if pos == 0 {
        return 0, 0, 0;
    }

    for pos > 2 {
        tmp := k + k0;
        k1 = k0;
        k0 = k;
        k = tmp;
        pos--;
    }
    return;
}

type FibResult struct {
    k, k0 uint64;    
}

func binFibJump(pos uint64) (uint64, uint64) {
    if pos == 2 {
        return 1, 1;
    }

    if pos == 1 {
        return 1, 0;
    }

    if pos == 0 {
        return 0, 0;
    }

    middle := pos / 2;
    f, f0 := binFibJump(middle);

    f1 := f - f0;
    unit := f + f0;
    
    v0 := f0 * unit + f1 * f;
    v  := f * unit + f * f0;
    
    if pos % 2 == 0 {
        return v, v0;
    }
    return v + v0, v;
}

func multiFibJump(pos uint64) uint64 {
    parts := runtime.NumCPU();

    if uint64(parts) > pos {
        parts = 1;
    }

    //interval := pos / uint64(parts)
    results := make([]FibResult, parts);
    outputs := make([]chan FibResult, parts);

    jmpSize := pos / uint64(parts);

    for i := range outputs {
        outputs[i] = make(chan FibResult, 10);
        go asyncFib(jmpSize, outputs[i]);
    }

    for i := range results {
        results[i] = <-outputs[i];
        close(outputs[i]);
    }

    var f, f0 uint64
    f, f0 = results[0].k, results[0].k0

    for i := 1; i < len(results); i++ {
        r := results[i]

        i := r.k
        i0 := r.k0

        i1 := i - i0;

        unit := f + f0;
        
        f0 = i0 * unit + i1 * f;
        f  = i * unit + i0 * f;
    }

    if exd := pos % uint64(parts); exd != 0 {
        i, i0, i1 := fibIndexer(exd);

        unit := f + f0;
        
        f0 = i0 * unit + i1 * f;
        f  = i * unit + i0 * f;
    }
    return f;
}


func asyncFib(pos uint64, out chan FibResult) {
    k, k0 := binFibJump(pos)
    out <- FibResult{k, k0};
}


func main() {
    makeTest();
}
