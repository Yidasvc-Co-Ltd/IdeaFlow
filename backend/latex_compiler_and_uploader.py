import subprocess
import shutil
import os

def compile_latex(latex_source, output_directory):
    # 生成临时的tex文件
    temp_tex_file = 'temp.tex'
    with open(temp_tex_file, 'w') as f:
        f.write(latex_source)

    try:
        # 调用pdflatex编译
        subprocess.run(['pdflatex', temp_tex_file])

        # 从生成的pdf文件中获取主文件名
        pdf_file = temp_tex_file.replace('.tex', '.pdf')

        # 按照指定规则重命名pdf文件
        new_pdf_name = 'new_name.pdf'  # 根据需要修改规则
        shutil.move(pdf_file, os.path.join(output_directory, new_pdf_name))
        print(f'PDF 文件已成功编译并重命名为: {new_pdf_name}')
    except Exception as e:
        print(f'编译失败: {e}')
    finally:
        # 删除临时tex文件
        os.remove(temp_tex_file)

if __name__ == '__main__':
    # LaTeX 源码
    latex_source_code = r'''
    \documentclass{article}
    \begin{document}
    Hello, LaTeX!
    \end{document}
    '''

    # 输出目录，可以是本地目录或者远程服务器目录（使用SSH传输）
    output_directory = '/root/data/upload/file/'

    # 调用编译函数
    compile_latex(latex_source_code, output_directory)
