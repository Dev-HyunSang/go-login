# Go-Login
> 본 프로젝트는 기본적인 로그인과 회원가입에 대해서 공부하고 개발하는 프로젝트입니다.  
> 본 프로젝트는 종료 되었습니다. 2021.08.21, 추후 JWT를 이용해서 다시 프로젝트를 열 예정입니다.  
> 프론트엔드에서의 처리가 아직 부족하여서 백엔드 개발을 다 완료 하였습니다.

본 프로젝트의 커밋 메시지 규칙은 [Conventinal Commit Messages](https://gist.github.com/qoomon/5dfcdf8eec66a051ecd85625518cfd13)를 따릅니다.

## 🚀 TODO:
### BackEnd
- [X] Render Handler | GET - 2021.07.19
    - `index.html` 완료 - 2021.07.19
    - `login.html` 완료 - 2021.07.19
    - `register.html` 완료 - 2021.07.19
- [X] `/register/new` POST: - 2021.07.19
    - [X] MySQL 연결 및 Table 구축
    - [X] `uuid` 패키지를 이용하여 `UUIDv4`로 개인별 ID 부여
    - [X] `bcrypt` 패키지를 이용하여 패스워드 암호화
- [ ] `/login` POST: 
    - [X] MySQL 연결  - 2021.07.22
    - [X] DB에서 UUID, Email, Password 가지고 오기 - 2021.07.26
    - [X] DB에서 가져온 정보를 암호화 된 패스워드와 사용자 입력 패스워드 대조  - 2021.07.27
    - [X] Session 정보를 MySQL에 입력 - 2021.08.01
        - 추가: 누가 로그인하였고 언제 했는지 추가함 / UUID / Email / TIME
    - [X] 로그인 완료시 JavaScript `alert` 띄우기
    - [ ] 로그인 완료시 Home으로 이동
        - 세션 처리 및 세션 확인 후 이동 시키기
  
## 기본 설정
### `.env` 설정하기
- `DB_ACCOUNT`: MySQL 계정
- `DB_HOST`: MySQL 주소
- `DB_PASSWORD`: MySQL 비밀번호
- `DB_NAME`: DataBase 이름

```env
DB_Connection_URL=
```

### DataBase Structure
회원가입시 사용되는 MySQL Table 
```sql
create  database go_login default character set utf8;
create table Users (
    ID  BINARY(36) primary key ,
    FirstName varchar(3) not null,
    LastName varchar(5) not null,
    Email varchar(320) not null,
    Password char(60) not null,
    CreatedAt Timestamp
);
```
로그인 후 세션 기록을 위하여 만든 MySQL Table
```SQL
create table autu_login (
    ID BINARY(36) primary key,
    Email varchar(320) not null,
    LogindAt timestamp
);
```
![FrontEnd Register](./images/Register.gif)
## 오류(고민) 해결기
### ID?
![error-01](./images/error-01.png)
사용자마다 다른 아이디를 생성하여서 고유 식별이 가능하도록 하여야 하는데 어떻게 해야할지 고민하던 중 코딩냄비에 질문을 해 본 결과 UUID를 사용하면 좋다고 해서 사용하게 됨.

### Session vs JWT 
![error-02](./images/error-02.png)
로그인 이후 관리를 어떻게 해야하는지 궁금해서 Discrod Gophers에 질문을 해 본 결과 보편적으로 사용하고 있는 방식은 Seesion 방식을 많이 사용하고 있다고 하여서 Seesion을 사용하기로 함.

### 참고하였던 자료
- DataBase
    - [예제로 배우는 Go 프로그래밍 - MySQL 사용 - 쿼리](http://golang.site/go/article/107-MySql-%EC%82%AC%EC%9A%A9---%EC%BF%BC%EB%A6%AC)
    - [[Go+MySQL] Go에서 MySQL 사용하기](https://soyoung-new-challenge.tistory.com/126)
- [Session](https://github.com/gorilla/sessions)
    - [Golang NewCookieStore Examples](https://golang.hotexamples.com/examples/github.com.gorilla.sessions/-/NewCookieStore/golang-newcookiestore-function-examples.html)
    - [session - 고릴라 세션을 사용하는 동안 golang의 세션 변수가 저장되지 않습니다](https://pythonq.com/so/session/457854)
    - [Go 언어 웹 프로그래밍 철저 입문 - 세션 관리](https://thebook.io/006806/ch09/03/01_01/)