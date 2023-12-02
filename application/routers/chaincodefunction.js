// const express = require("express")
// const fs = require("fs")
// const path = require("path")
// const { Gateway, Wallets } = require("fabric-network");

// const router = express.Router()

// router.post("/registeracm", async (req,res)=>{
//     // 사용자가 입력한 숙소 정보 읽어오기
//     // 쿠키에서 userid 읽어오기
//     const {acmlocate,acmadr,acmname,acmroom,acmprice} = req.body
//     const userCookie = req.cookies[`USER`]
//     userData = JSON.parse(userCookie) //userData.userid
    

// async function
//     // 체인코드 RegisterAccommodation 함수 호출 
//     // 호출한 결과를 웹브라우저에게 전송(res)
//     // res.render("registeracm")
//     try {
//         // load the network configuration
//         const ccpPath = path.resolve(__dirname, "..", "ccp", "connection-org1.json");
//         let ccp = JSON.parse(fs.readFileSync(ccpPath, "utf8"));

//         // Create a new file system based wallet for managing identities.
//         const walletPath = path.join(process.cwd(), "wallet");
//         const wallet = await Wallets.newFileSystemWallet(walletPath);
//         console.log(`Wallet path: ${walletPath}`);

//         // Check to see if we've already enrolled the user.
//         const identity = await wallet.get("appUser");
//         if (!identity) {
//             console.log(
//                 'An identity for the user "appUser" does not exist in the wallet'
//             );
//             console.log("Run the registerUser.js application before retrying");
//             return;
//         }

//         // Create a new gateway for connecting to our peer node.
//         const gateway = new Gateway();
//         await gateway.connect(ccp, {
//             wallet,
//             identity: "appUser",
//             discovery: { enabled: true, asLocalhost: true },
//         });

//         // Get the network (channel) our contract is deployed to.
//         const network = await gateway.getNetwork("mychannel");

//         // Get the contract from the network.
//         const contract = network.getContract("myBooking");

//         // Submit the specified transaction.
//         // createCar transaction - requires 5 argument, ex: ('createCar', 'CAR12', 'Honda', 'Accord', 'Black', 'Tom')
//         // changeCarOwner transaction - requires 2 args , ex: ('changeCarOwner', 'CAR12', 'Dave')
//         await contract.submitTransaction(
//             "RegisterAccommodation",
//             userData.userid,
//             userData.auth,
//             acmname,
//             acmroom,
//             acmprice,
//             acmlocate,
//             acmadr
//         );
//         console.log("Transaction has been submitted");
 
//         // Disconnect from the gateway.
//         await gateway.disconnect();
//     } catch (error) {
//         console.error(`Failed to submit transaction: ${error}`);
//         // res_str = `{"resultcode":"failed", "msg":"자산생성 실패"}`
//         // res.json(JSON.parse(res_str))
//         res.write("<script>alert('자산생성 실패')</script>")
//         return
    
//     }

//     // 블록체인 연동후 체인코드 호출이 정상적으로 잘 처리된경우
//     // 그 실행 결과를 받아와서 response를 만들어 웹 브라우저에 전달
//     // res_str = `{"resultcode":"success", "msg":"자산생성이 정상적으로 완료되었습니다."}`
//     // res.json(JSON.parse(res_str))

//     res.redirect("/")

// })


// router.get('/updateacm',(req,res)=>{
    
//     res.render("updateacm")
// })

// router.get('/readacm',(req,res)=>{
    
//     res.render("readacm")
// })
// module.exports = router //외부에서 export해서 쓸 수 있도록 함


// //  invoke 함수 호출부분 -함수로 만들기

// func (){
//     if (fn_name =="RegisterAccommodation"){
//         await contract.submitTransacit(
//             "RegisterAccommodation",
//             args[0],
//             args[1],
//             args[2]
//         );
    
//     }

// }