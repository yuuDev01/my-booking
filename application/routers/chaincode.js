const express = require("express")
const fs = require("fs")
const path = require("path")
const { Gateway, Wallets } = require("fabric-network");

const router = express.Router()

router.post("/registeracm", async (req,res)=>{
    // 사용자가 입력한 숙소 정보 읽어오기
    // 쿠키에서 userid 읽어오기
    const {acmlocate,acmadr,acmname,acmroom,acmprice} = req.body
    const userCookie = req.cookies[`USER`]
    userData = JSON.parse(userCookie) //userData.userid
    
    
    // 체인코드 RegisterAccommodation 함수 호출 
    // 호출한 결과를 웹브라우저에게 전송(res)
    // res.render("registeracm")
    try {
        // load the network configuration
        const ccpPath = path.resolve(__dirname, "..", "ccp", "connection-org1.json");
        let ccp = JSON.parse(fs.readFileSync(ccpPath, "utf8"));

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), "wallet");
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const identity = await wallet.get("appUser");
        if (!identity) {
            console.log(
                'An identity for the user "appUser" does not exist in the wallet'
            );
            console.log("Run the registerUser.js application before retrying");
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccp, {
            wallet,
            identity: "appUser",
            discovery: { enabled: true, asLocalhost: true },
        });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork("mychannel");

        // Get the contract from the network.
        const contract = network.getContract("myBooking");

        // Submit the specified transaction.
        // createCar transaction - requires 5 argument, ex: ('createCar', 'CAR12', 'Honda', 'Accord', 'Black', 'Tom')
        // changeCarOwner transaction - requires 2 args , ex: ('changeCarOwner', 'CAR12', 'Dave')
        await contract.submitTransaction(
            "RegisterAccommodation",
            userData.userid,
            userData.userauth,
            acmname,
            acmroom,
            acmprice,
            acmlocate,
            acmadr
        );
        console.log("Transaction has been submitted");
        console.log(userData.userid,
            userData.userauth,
            acmname,
            acmroom,
            acmprice,
            acmlocate,
            acmadr)
        // Disconnect from the gateway.
        await gateway.disconnect();
    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        // res_str = `{"resultcode":"failed", "msg":"자산생성 실패"}`
        // res.json(JSON.parse(res_str))
        res.write("<script>alert('자산생성 실패')</script>")
        return
    
    }

    // 블록체인 연동후 체인코드 호출이 정상적으로 잘 처리된경우
    // 그 실행 결과를 받아와서 response를 만들어 웹 브라우저에 전달
    // res_str = `{"resultcode":"success", "msg":"자산생성이 정상적으로 완료되었습니다."}`
    // res.json(JSON.parse(res_str))

    res.redirect("/")

})

// 내가 등록한 숙소 조회
router.get('/myacm',  async (req,res)=>{
    const userCookie = req.cookies[`USER`]
    userData = JSON.parse(userCookie) //userData.userid
    console.log(userData.userid)
    try {
        // load the network configuration
        const ccpPath = path.resolve(__dirname,"..", "ccp", "connection-org1.json");
        let ccp = JSON.parse(fs.readFileSync(ccpPath, "utf8"));

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), "wallet");
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const identity = await wallet.get("appUser");
        if (!identity) {
            console.log(
                'An identity for the user "appUser" does not exist in the wallet'
            );
            console.log("Run the registerUser.js application before retrying");
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccp, {
            wallet,
            identity: "appUser",
            discovery: { enabled: true, asLocalhost: true },
        });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork("mychannel");

        // Get the contract from the network.
        const contract = network.getContract("myBooking");

        // Submit the specified transaction.
        // createCar transaction - requires 5 argument, ex: ('createCar', 'CAR12', 'Honda', 'Accord', 'Black', 'Tom')
        // changeCarOwner transaction - requires 2 args , ex: ('changeCarOwner', 'CAR12', 'Dave')
        result = await contract.evaluateTransaction(
            "GetMyAccommodations",
            userData.userid
        );
        console.log("Transaction has been submitted");
        res_str = `{"resultcode":"success", "msg": ${result}}`
        res.json(JSON.parse(res_str))
        

        // Disconnect from the gateway.
        await gateway.disconnect();
    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        res_str = `{"resultcode":"failed", "msg":"자산조회 실패"}`
        res.json(JSON.parse(res_str))
    
    }

})
router.get('/updateacm',(req,res)=>{
    
    res.render("updateacm")
})

// 수정하기.UpdateAccommodation(uid string, aid string, price int, available bool)
router.post('/updateacm', async (req, res)=>{
    const {acmid, acmprice, acmavailable} = req.body
    const userCookie = req.cookies[`USER`]
    userData = JSON.parse(userCookie) //userData.userid

    try {
        // load the network configuration
        const ccpPath = path.resolve(__dirname, "..", "ccp", "connection-org1.json");
        let ccp = JSON.parse(fs.readFileSync(ccpPath, "utf8"));

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), "wallet");
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const identity = await wallet.get("appUser");
        if (!identity) {
            console.log(
                'An identity for the user "appUser" does not exist in the wallet'
            );
            console.log("Run the registerUser.js application before retrying");
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccp, {
            wallet,
            identity: "appUser",
            discovery: { enabled: true, asLocalhost: true },
        });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork("mychannel");

        // Get the contract from the network.
        const contract = network.getContract("myBooking");

        // Submit the specified transaction.
        // createCar transaction - requires 5 argument, ex: ('createCar', 'CAR12', 'Honda', 'Accord', 'Black', 'Tom')
        // changeCarOwner transaction - requires 2 args , ex: ('changeCarOwner', 'CAR12', 'Dave')
        await contract.submitTransaction(
            "UpdateAccommodation",
            userData.userid,
            acmid,
            acmprice,
            acmavailable
        );
        console.log("Transaction has been submitted");
        console.log(
            userData.userid,
            acmid,
            acmprice,
            acmavailable
            )
        // Disconnect from the gateway.
        await gateway.disconnect();
    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        res_str = `{"resultcode":"failed", "msg":"수정 실패"}`
        res.json(JSON.parse(res_str))
        return
    
    }

    // 블록체인 연동후 체인코드 호출이 정상적으로 잘 처리된경우
    // 그 실행 결과를 받아와서 response를 만들어 웹 브라우저에 전달
    res_str = `{"resultcode":"success", "msg":"수정 완료되었습니다."}`
    res.json(JSON.parse(res_str))

    // res.redirect("/")
})

// 숙소 삭제 (invoke)
router.get('/deleteacm', async (req, res)=>{
    const acmid = req.query.acmid
    const userCookie = req.cookies[`USER`]
    userData = JSON.parse(userCookie) //userData.userid

    console.log(acmid)
    try {
        // load the network configuration
        const ccpPath = path.resolve(__dirname, "..", "ccp", "connection-org1.json");
        let ccp = JSON.parse(fs.readFileSync(ccpPath, "utf8"));

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), "wallet");
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const identity = await wallet.get("appUser");
        if (!identity) {
            console.log(
                'An identity for the user "appUser" does not exist in the wallet'
            );
            console.log("Run the registerUser.js application before retrying");
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccp, {
            wallet,
            identity: "appUser",
            discovery: { enabled: true, asLocalhost: true },
        });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork("mychannel");

        // Get the contract from the network.
        const contract = network.getContract("myBooking");

        // Submit the specified transaction.
        // createCar transaction - requires 5 argument, ex: ('createCar', 'CAR12', 'Honda', 'Accord', 'Black', 'Tom')
        // changeCarOwner transaction - requires 2 args , ex: ('changeCarOwner', 'CAR12', 'Dave')
        await contract.submitTransaction(
            "DeleteAccommodation",
            userData.userid,
            acmid,
        );
        console.log("Transaction has been submitted");
        console.log(
            userData.userid,
            acmid,

            )
        // Disconnect from the gateway.
        await gateway.disconnect();
    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        res_str = `{"resultcode":"failed", "msg":"수정 삭제 실패"}`
        res.json(JSON.parse(res_str))
        // res.write("<script>alert('수정 실패')</script>")
        return
    
    }

    // 블록체인 연동후 체인코드 호출이 정상적으로 잘 처리된경우
    // 그 실행 결과를 받아와서 response를 만들어 웹 브라우저에 전달
    res_str = `{"resultcode":"success", "msg":"수정 완료되었습니다."}`
    res.json(JSON.parse(res_str))

    // res.redirect("/")
})

module.exports = router //외부에서 export해서 쓸 수 있도록 함