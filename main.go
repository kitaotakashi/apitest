//https://qiita.com/katekichi/items/d94e078b376151858ca4
//https://qiita.com/suin/items/f32fa82d6c35a34e8d16
//https://qiita.com/hirokiseiya/items/8a2c9d6f6e0fb5db329f
package main

import (
    "fmt"
    "net/http"
)

func main() {
    // buildフォルダを公開
    //http.Handle("/admin/", http.StripPrefix("/admin/", http.FileServer(http.Dir("./build"))))

    // 前述のhomepageを記述してないと下の記述でしか動かない。他のパスだとバグる
    http.Handle("/", http.FileServer(http.Dir("./robot_front/build")))

    fmt.Println("Server Started Port 443 v2")
    http.ListenAndServe(":443", nil)
}
