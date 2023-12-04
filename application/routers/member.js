const express = require("express")
const fs = require("fs").promises


const router = express.Router()
// 유저
//숙소 조회
router.get("/acmlist",(req,res)=>{
    res.render("acmlist")
})



// 숙소 예약
router.get("/bookingcreate",(req,res)=>{
    res.render("bookingcreate")
})

// 내 예약 조회
router.get("/bookingread",(req,res)=>{
    res.render("bookingread")
})

// 예약수정
router.get("/bookingupdate",(req,res)=>{
    res.render("bookingupdate")
})

// 예약취소
router.get("/bookingdelete",(req,res)=>{
    res.render("bookingdelete")
})

// 숙소 조회
router.get('/acmlist',(req,res)=>{

})

module.exports = router 