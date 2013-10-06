JIT
---

A `libjit` wrapper for Golang.

(c) 2013 Chanwit Kaewkasi

Suranaree University of Technology, Thailand

How to install
--------------

You need to install `libjit` firstly. Please note that it requires `automake 1.11.6` or later to build `libjit`. So, check your version of `automake` before proceed.

    git clone git://git.savannah.gnu.org/libjit.git
    cd libjit
    ./auto_gen.sh
    ./configure --prefix=/usr
    make
    sudo make install

Then, install the library by: `go get github.com/chanwit/jit`.

There are some examples under `jit/examples`. You can try running `go run t1.go` and see the result.

Enjoy jitting.