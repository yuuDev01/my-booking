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

// 미들웨어 등록부분


app.use(express.static(path.join(__dirname, "views")));
app.use(bodyParser.json())
app.use(bodyParser.urlencoded({extended:false}))

app.set('views', __dirname + '/views');
app.set('view engine', 'ejs');

app.use("/user", userRouter)
app.get("/", (req,res)=>{
    //로그인완료
    res.render("index")
    //
})

app.listen(PORT, HOST)
console.log(`Running on http://${HOST}:${PORT}`)