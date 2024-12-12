import os
from flask import Flask, request, render_template, send_file
from main import judge_zhihu_type

app = Flask(__name__)

@app.route("/", methods=["GET", "POST"])
def index():
    if request.method == "POST":
        cookies = request.form["cookies"]
        urls = request.form["url"].strip().split('\n')
        
        # 创建下载目录
        download_dir = "downloads"
        os.makedirs(download_dir, exist_ok=True)
        
        # 保存当前工作目录
        old_cwd = os.getcwd()
        os.chdir(download_dir)
        
        # 记录失败的URL到日志文件
        error_log_path = os.path.join(old_cwd, 'error.log')
        
        # 处理每个URL
        for url in urls:
            url = url.strip()
            if url:  # 跳过空行
                try:
                    judge_zhihu_type(url, cookies, None)
                except Exception as e:
                    error_msg = f"处理URL时出错: {url}, 错误: {str(e)}\n"
                    print(error_msg)  # 控制台输出
                    # 写入错误日志
                    with open(error_log_path, 'a', encoding='utf-8') as f:
                        f.write(error_msg)
                    continue
        
        # 恢复工作目录
        os.chdir(old_cwd)
        
        return {"message": "文件已保存到downloads目录"}

    return render_template("index.html")

@app.route("/get-cookies")
def get_cookies():
    return render_template("howToGetCookies.html")

if __name__ == "__main__":
    app.run(host="0.0.0.0", debug=False)
