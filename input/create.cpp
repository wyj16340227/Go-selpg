/*
 * 	生成一个测试文件
 * 		文件共有500行，每行格式为： line + 行号
 * 			行号从1开始，直到500
 * 				每10行有一个换页符'\f'将行号隔开，因此共有50页
 * 					输出在data.txt文件中
 * 					*/
#include <iostream>
#include <fstream>
#include <string>

using namespace std;

int main () {
	ofstream fout;
	fout.open("data.com");
	string line = "line ";
	string page = "\f";
	for (int i = 0; i < 50; i++) {
		for (int j = 0; j < 10; j++) {
			fout << line << i * 50 + j + 1 << endl;
		}
		fout << page;
	}
}
