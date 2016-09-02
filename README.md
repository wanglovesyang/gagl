# gaql

GAQL is short for Golang-Augmented Querying Language, which is a SQL-like querying language conducted on distributed data set (eg. HDFS). Some enhanced features beyond convertional SQL querying are integrated in the language, enabling the users to implement their own mapping/reducing/joining logic in the SQL query.  

The GAQL project includes following components:
1. A enhanced SQL-like language with syntax analyzer.
2. A cross-compiler which links GAQL scripts with *.go files.
3. An instant deploy & execution platform for hadoop yarn.

This project is inspired by the SCOPE language of Microsoft, which allows integration of use-defined moufules in C#. Besides go language and open source, there are several other differences: