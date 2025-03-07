package main

import (
    "fmt"
    "time"
    "runtime"
    "math/big"
)

func bigI(a int64) *big.Int {
    return big.NewInt(a);
}

func sum(a, b *big.Int) *big.Int {
    var z big.Int;
    return z.Add(a, b);
}

func sub(a, b *big.Int) *big.Int {
    var z big.Int;
    return z.Sub(a, b);
}

func mul(a, b *big.Int) *big.Int {
    var m big.Int;
    return m.Mul(b, a);
}

func fib(pos int64) *big.Int {
    if pos <= 1 {
        return big.NewInt(pos);
    }
    return sum(fib(pos - 1), fib(pos - 2));
}

func fibJump(pos int64) *big.Int {
    if pos <= 1 {
        return bigI(pos);
    }

    const psize = 20;
    
    i, i0, i1 := fibIndexer(psize);
    f, f0 := i, i0;

    for iter := int64(1); iter < pos / psize; iter++ {
        unit := sum(f, f0);

        f0 = sum(mul(i0, unit), mul(i1, f));
        f  = sum(mul(i, unit),  mul(i0, f));
    }

    if remain := pos % psize; remain != 0 {
        i, i0, i1 := fibIndexer(remain); 
        unit := sum(f, f0);

        f0 = sum(mul(i0, unit), mul(i1, f));
        f  = sum(mul(i, unit),  mul(i0, f));
    }
    return f;
}

func fibIndexer(pos int64) (k, k0, k1 *big.Int) {
    k, k0, k1 = bigI(1), bigI(1), bigI(0);

    if pos == 2 {
        return;
    }

    if pos == 1 {
        return bigI(1), bigI(0), bigI(0);
    }

    if pos == 0 {
        return bigI(0), bigI(0), bigI(0);
    }

    for pos > 2 {
        tmp := sum(k, k0);
        k1 = k0;
        k0 = k;
        k = tmp;
        pos--;
    }
    return;
}

func binFibJump(pos int64) (*big.Int, *big.Int) {
    if pos == 2 {
        return bigI(1), bigI(1);
    }

    if pos == 1 {
        return bigI(1), bigI(0);
    }

    if pos == 0 {
        return bigI(0), bigI(0);
    }

    middle := pos / 2;
    f, f0 := binFibJump(middle);

    f1 := sub(f, f0);
    unit := sum(f, f0);
    
    v0 := sum(mul(f0, unit), mul(f1, f));
    v  := sum(mul(f, unit), mul(f, f0));
    
    if pos % 2 == 0 {
        return v, v0;
    }
    return sum(v, v0), v;
}

type FibResult struct {
    k, k0 *big.Int;
};

func multiFibJump(pos int64) *big.Int {
    parts := runtime.NumCPU();

    if int64(parts) > pos {
        parts = 1;
    }

    //interval := pos / uint64(parts)
    results := make([]FibResult, parts);
    outputs := make([]chan FibResult, parts);

    jmpSize := pos / int64(parts);

    for i := range outputs {
        outputs[i] = make(chan FibResult, 10);
        go asyncFib(jmpSize, outputs[i]);
    }

    for i := range results {
        results[i] = <-outputs[i];
        close(outputs[i]);
    }

    var f, f0 *big.Int
    f, f0 = results[0].k, results[0].k0

    for i := 1; i < len(results); i++ {
        r := results[i]

        i := r.k
        i0 := r.k0

        i1 := sub(i, i0);

        unit := sum(f, f0);
        
        f0 = sum(mul(i0, unit), mul(i1, f));
        f  = sum(mul(i, unit),  mul(i0, f));
    }

    if exd := pos % int64(parts); exd != 0 {
        i, i0, i1 := fibIndexer(exd);

        unit := sum(f, f0);
        
        f0 = sum(mul(i0, unit), mul(i1, f));
        f  = sum(mul(i, unit),  mul(i0, f));
    }
    return f;
}

func fibSeq(cache []*big.Int) {
    if len(cache) < 2 {
        return;
    }

    if cache[0] == nil {
        cache[0] = bigI(0);
        cache[1] = bigI(1);
    }

    for i := 2; i < len(cache); i++ {
        cache[i] = sum(cache[i - 1], cache[i - 2]);
    }
}

func mFibSeq(cache []*big.Int) {
    parts := runtime.NumCPU();

    if len(cache) < parts {
        fibSeq(cache);
        return;
    }

    channels := make([]chan bool, parts);
    for i := range channels {
        channels[i] = make(chan bool, 10);
    }

    interval := len(cache) / parts;

    for i := 0; i < parts; i++ {
        pos := i * interval;
        var f, f0 *big.Int 

        if i != 0 {
            f, f0 = binFibJump(int64(pos));
            tmp := f0;
            f0 = f;
            f = sum(f, tmp);

        } else {
            f, f0 = bigI(1), bigI(0);
        }

        slice := cache[pos:pos+interval];
        slice[1] = f;
        slice[0] = f0;

        fmt.Printf("Slicing %d, %d\n", pos, pos+interval);
        fmt.Println("s[0] = ", slice[0]);
        fmt.Println("s[1] = ", slice[1]);
        go asyncFibSeq(slice, channels[i]);
    }

    for _, c := range channels {
        <-c;
        close(c);
    }

    if ext := len(cache) % parts; ext != 0 {
        slice := cache[parts*interval-2:];
        fibSeq(slice);
    }
}

func asyncFibSeq(cache []*big.Int, ended chan bool) {
    fibSeq(cache);
    ended <- true;
}

func asyncFib(pos int64, out chan FibResult) {
    k, k0 := binFibJump(pos)
    out <- FibResult{k, k0};
}

func makeTest() {
    titles := []string{ "Fib seq", "Fib multi seq" };
    fn := []func([]*big.Int){ fibSeq, mFibSeq };

    for i := range titles {
        fmt.Println("---", titles[i]);
        cache := make([]*big.Int, 100);

        start := time.Now();
        fn[i](cache);
        end := time.Now().Sub(start);

        for i, c := range cache {
            fmt.Printf("fib(%d) = %d\n", i, c);
        }
        fmt.Println("Time:", end);
    }
}

func main() {
    makeTest();
}
