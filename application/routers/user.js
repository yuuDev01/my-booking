const express = require("express")
const fs = require("fs").promises


const router = express.Router()

router.post("/login", async (req,res)=>{
    console.log("in login")
    const {userid, userpw} = req.body
    const user = await fetchUsers(userid)
    if(!user){
        res.send("존재하지 않는 사용자")
        return
    }
    if(userpw !== user.userpw){
        res.send("비밀번호 불일치")
        return
    }
    // res.send("로그인 성공"), 쿠키 생성
    cookie_str = `{"userid":"${user.userid}", "username":"${user.username}", "userauth":"${user.userauth}" }`
    res.cookie("USER", cookie_str)

    res.redirect("/")
})


router.get('/logout',(req,res)=>{
    res.clearCookie("USER")
    res.redirect("/")
})

router.post("/register", (req,res)=>{
    res.render("signup")
})

router.post("/signup", async (req,res)=>{
    // 회원가입 수행 
    const {userid, userpw, username, userauth} = req.body
    const newUser = {userid, userpw, username, userauth}

    const exist = await fetchUsers(userid)
    if(exist){
        res.send("이미 존재하는 사용자 입니다.")
        return
    }
    await createUser(newUser)
    res.redirect("/")
})

// 모든 데이터 가져와서 구조체로 리턴
async function fetchAllUsers(){
    const userdata = await fs.readFile("./public/userdata/user.json")
    const users =  JSON.parse(userdata.toString())
    return users
}

// 회원가입 
async function createUser(newUser){
    const users = await fetchAllUsers()
    users.push(newUser) //변수 user에 새로운 유저데이터 넣기 
    await fs.writeFile("./public/userdata/user.json", JSON.stringify(users)) //유저데이터 넣은 user.json파일생성, 기존 파일이 있으면 대체 

}

// 이미 존재하는 아이디인지 확인
async function fetchUsers(userid){
    const users = await fetchAllUsers()
    const user = users.find((users)=>users.userid === userid)
    return user
}


module.exports = router //외부에서 export해서 쓸 수 있도록 함

