const express = require("express")
const router = express.Router()

router.get("/registeracm", async (req,res)=>{
    res.render("registeracm")
})

router.get('/myacm/:userid', (req,res)=>{
    const myid = req.params.userid
    res.render("myacm", { myid: myid })
})

router.get('/updateacm',(req,res)=>{
    
    res.render("updateacm")
})



module.exports = router //외부에서 export해서 쓸 수 있도록 함

