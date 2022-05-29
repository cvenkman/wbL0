./publisher -ch=boo -data=testData/model1.json
echo "order_uid = 1abc"
./publisher -ch=boo -data=testData/model2.json
echo "order_uid = 2"
# ./publisher -ch=boo -data=testData/model3.json
# try publish files with errors
./publisher -ch=boo -data=testData/modelError.json
./publisher -ch=boo -data=testData/err.sh
./publisher -ch=boo -data=testData/emptyError.json
# try publish model that already exists
./publisher -ch=boo -data=testData/model2.json