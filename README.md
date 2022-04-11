## 근태관리프로그램
  
    제작기간 : 약 1달
    제작언어 : Go v1.10
    구동환경 : - Server [ linux ] ( module base 'ECHO' )
              - Mariadb 10.x [ linux ]
              - Cleint [ win ] ( module base 'LXW' )
    구동방식 : Direct Rest API 
          
## 사용방법

    각자의 직원은 직장 출근시 클라이언트 프로그램을 실행하여 출근체크 버튼을 누릅니다.
    클라이언트는 Windows gui 형태의 exe 파일이며 서버에 제공하는 Rest URL에 출근 정보를 태워서 보낸다.
    서버에서는 해당 정보를 DB에 저장하여 근태기록을 볼 수 있도록 한다.

