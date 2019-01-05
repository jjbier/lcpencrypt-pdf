
## [lcpencrypt-pdf]  

A command line utility for PDF content encryption. This utility can be included in any processing pipeline. 

* takes one unprotected PDF file as input and generates an encrypted file as output.
* notifies the License server of the generation of an encrypted file.


## Example
-input /workspace/go/libros/103_1415970848_5466002100a8e_5466d2fc4bc8a.pdf -output /workspace/lcpconfig/lcpfiles/encrypted/103_1415970848_5466002100a8e_5466d2fc4bc8a.pdf -lcpsv http://127.0.0.1:8989 -login admin -password admin 
./lcpencrypt-pdf -input /workspace/go/libros/103_1415970848_5466002100a8e_5466d2fc4bc8a.pdf -output /workspace/lcpconfig/lcpfiles/encrypted/103_1415970848_5466002100a8e_5466d2fc4bc8a.pdf -lcpsv http://127.0.0.1:8989 -login admin -password admin 