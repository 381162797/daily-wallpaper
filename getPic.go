package main


import (
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "time"
    "io"
    "encoding/json"
    "os/user"
    "os/exec"
    "path/filepath"
)


const (
    baseUrl = "https://cn.bing.com"
    apiUrl  = baseUrl + "/HPImageArchive.aspx?format=js&idx=0&n=1"
    STMENU = `C:\ProgramData\Microsoft\Windows\Start Menu\Programs\StartUp\wallpaper.exe`
)

func getImgUrl() (string,error) { 
    type Images struct{
        Images []map[string]interface{} `json:"images"`
        Tooltips map[string]string `json:"tooltips"`
    }
    var myImages Images
    resp, err := http.Get(apiUrl)
	if err != nil {
		return "",err
	}
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
		return "",err
	}
    err = json.Unmarshal(body,&myImages)
    if err != nil{
        fmt.Println("json" ,err)
    }
    imgUrl := myImages.Images[0]["url"]
    return imgUrl.(string),nil
}

func saveImg(url string) error{
    res, err := http.Get(baseUrl + url)
	if err != nil {
		return err
    }
    cUser, err := user.Current()
	if err != nil {
		return err
    }
    desktop := cUser.HomeDir + "/Desktop/images"
    name := "/" + time.Now().Format("2006-01-02")
    os.Mkdir(desktop,0633)
	f, err := os.Create(desktop+name+".jpg")
	if err != nil {
		return err
	}
    io.Copy(f, res.Body)
    return nil
}

func start(){
    file, _ := exec.LookPath(os.Args[0])
    path, _ := filepath.Abs(file)
    fmt.Println(path)
    if path != STMENU {
      //拷贝进启动目录
      reader,err := ioutil.ReadFile(path)
      if err != nil {
          fmt.Println(err)
      }else{
        err = ioutil.WriteFile(STMENU,reader,0633)
        if err != nil {
            fmt.Println(err)
        }
      }
    }
}


func main() {
    var keys int
    start()
    for {
        _ , err := http.Get("http://www.baidu.com")
        if err != nil {
            keys ++
            time.Sleep(time.Minute)
        }else{
            break
        }
        if keys >= 10 {
            return
        }
        
    }

    url,err := getImgUrl()
    if err != nil {
        fmt.Println("getImgUrl error",err)
    }else{
        err := saveImg(url)
        if err != nil {
            fmt.Println("save img error",err)
        }
    }

}

