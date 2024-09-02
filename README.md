# edu-container
Simple container runtime for educational  purposes

inspired by https://youtu.be/sK5i-N34im8?si=QQYBdm3RCn8uxa3T and https://www.youtube.com/watch?v=8fi7uSYlOdc


# how to run
1. clone the code onto a linux system
   - ```git clone https://github.com/maexoman/edu-container.git```
   - ```cd edu-container``` 
2. get a copy of a basic linux filesystem
   - ```sudo docker export $(sudo docker create alpine) --output="apline.tar"```
   - ```mkdir layer1```
   - ```sudo tar -xf apline.tar -C ./layer1```
   - ```sudo rm apline.tar```
3. set the path to the layer in the layers splice (please use an absolute path)
4. change the root path to your home direcotry (please use an absolute path)
5. run ```sudo go run . run /bin/sh```
6. have fun :)
