package main

/*===================import=====================*/

import (
	"bufio"
	"fmt"
	flag "github.com/spf13/pflag"		//pflag包
	"io"
	"os"
	"os/exec"
	"strings"
)

/*====================type===================*/

type selpg_Args struct {
	start int				//起始页码
	end int					//终止页码
	length int				//行数
	readType string			//‘l’代表按照行数计算页；‘f’为按照换页符计算页
	dest string				//定向位置
	inputType string		//输入方式

}

/*====================function=================*/

func main() {
	sa := new(selpg_Args)

	//参数绑定变量
	flag.IntVar(&sa.start, "s", 0, "the start Page")				//开始页码，默认为0
	flag.IntVar(&sa.end, "e", 0, "the end Page")					//结束页码，默认为0
	flag.IntVar(&sa.length, "l", 72, "the length of the page")		//每页行数，默认为72行每页
	flag.StringVar(&sa.dest, "d", "", "the destiny of printing")	//输出位置，默认为空字符

	//查找 f
	isF := flag.Bool("f", false, "")
	flag.Parse()

	//如果输入f，按照f并取-1；否则按照 l
	if *isF {
		sa.readType = "f"
		sa.length = -1
	} else {
		sa.readType = "l"
	}

	//如果使用了文件输入，将方式置为文件名
	sa.inputType = ""
	if flag.NArg() == 1 {
		sa.inputType = flag.Arg(0)
	}


/*====================check=================*/
	//检查剩余参数数量
	if narg := flag.NArg(); narg != 1 && flag.NArg() != 0 {
		usage()
		os.Exit(1)
	}
	//检查起始终止页
	if sa.start > sa.end || sa.start < 1 {
		usage()
		os.Exit(1)
	}
	//检查l f 是否同时出现
	if sa.readType == "f" && sa.length != -1 {
		usage()
		os.Exit(1)
	}

	run(*sa)
}

//执行指令
func run(sa selpg_Args) {
	//初始化
	fin := os.Stdin					//输入
	fout := os.Stdout				//输出
	currentLine := 0				//当前行
	currentPage := 1				//当前页
	var inpipe io.WriteCloser		//管道
	var err error					//错误

	//判断输入方式
	if sa.inputType != "" {
		fin, err = os.Open(sa.inputType)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Can't find input file \"%s\"!\n", sa.inputType)
			//fmt.Println(err)
			usage()
			os.Exit(1)
		}
		defer fin.Close()			//全部结束了再关闭
	}

	//确定输出到文件或者输出到屏幕
	//通过用管道接通grep模拟打印机测试，结果输出到屏幕
	if sa.dest != "" {
		cmd := exec.Command("grep", "-nf", "keyword")
		inpipe, err = cmd.StdinPipe()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer inpipe.Close() //最后执行
		cmd.Stdout = fout
		cmd.Start()
	}

	//分页方式
	//设置行数
	if sa.readType == "l" {
		//按照行读取
		line := bufio.NewScanner(fin)
		for line.Scan() {
			if currentPage >= sa.start && currentPage <= sa.end {
				//输出到窗口
				fout.Write([]byte(line.Text() + "\n"))
				if sa.dest != "" {
					//定向到文件管道
					inpipe.Write([]byte(line.Text() + "\n"))
				}
			}
			currentLine++
			//翻页
			if currentLine == sa.length {
				currentPage++
				currentLine = 0
			}
		}
	} else {
		//用换行符 '\f'分页 
		rd := bufio.NewReader(fin)
		for {
			page, ferr := rd.ReadString('\f')
			if ferr != nil || ferr == io.EOF {
				if ferr == io.EOF {
					if currentPage >= sa.start && currentPage <= sa.end {
						fmt.Fprintf(fout, "%s", page)
					}
				}
				break
			}
			//'\f'翻页
			page = strings.Replace(page, "\f", "", -1)
			currentPage++
			if currentPage >= sa.start && currentPage <= sa.end {
				fmt.Fprintf(fout, "%s", page)
			}
		}
	}
	//当输出完成后，比较输出的页数与期望输出的数量
	if currentPage < sa.end {
		fmt.Fprintf(os.Stderr, "./selpg: end (%d) greater than total pages (%d), less output than expected\n", sa.end, currentPage)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "\nUSAGE: ./selpg [--s start] [--e end] [--l lines | --f ] [ --d dest ] [ in_filename ]\n")
	fmt.Fprintf(os.Stderr, "\n selpg --s start    : start page")
	fmt.Fprintf(os.Stderr, "\n selpg --e end      : end page")
	fmt.Fprintf(os.Stderr, "\n selpg --l lines    : lines/page")
	fmt.Fprintf(os.Stderr, "\n selpg --f          : check page with '\\f'")
	fmt.Fprintf(os.Stderr, "\n selpg --d dest     : pipe destination\n")
}
