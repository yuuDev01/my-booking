const express = require("express")
const bodyParser = require("body-parser");
const path = require("path")
const cookieParser = require("cookie-parser")
const fs = require("fs");
const { Gateway, Wallets } = require("fabric-network");

const app = express()

const PORT = 3000;
const HOST = "0.0.0.0";

// 라우터 등록 부분
var userRouter = require("./routers/user")
var adminRouter = require("./routers/admin")
var chaincodeRouter = require("./routers/chaincode")

// 미들웨어 등록부분
app.use(express.static(path.join(__dirname, "views")));
app.use(bodyParser.json())
app.use(bodyParser.urlencoded({extended:false}))
app.use(cookieParser());

app.set('views', __dirname + '/views');
app.set('view engine', 'ejs');

app.use("/user", userRouter)
app.use("/admin", adminRouter)
app.use("/chaincode", chaincodeRouter)

app.get("/", (req,res)=>{
 const userCookie = req.cookies[`USER`]
 
 //로그인완료
 console.log(userCookie)
 
 //로그인이 안되어있을 경우
 if(!userCookie){
     res.render("index")
}else{
     //로그인 성공
     // 유저 쿠키 객체 생성
     userData = JSON.parse(userCookie)
     console.log(userCookie)
     //한페이지에 관리자, 사용자에 따라 다른 화면 나오게 하기
    //  if(userData.auth === "admin"){
    //      res.render("acmlist", {userclass: "admin", username:userData.username})
    // }else{
    //     res.render("acmlist", {userclass: "user", username:userData.username})   
    //     }
    
    // 관리자 , 사용자 화면 따로 생성하여 연결
    if(userData.userauth === "admin"){
             res.render("adminidx", {username:userData.username, userid:userData.userid})
        }else{
            res.render("useridx", {username:userData.username})   
            }
}

    
    // console.log(auth)
    

})

app.listen(PORT, HOST)
console.log(`Running on http://${HOST}:${PORT}`)