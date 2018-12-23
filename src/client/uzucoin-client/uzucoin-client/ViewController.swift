//
//  ViewController.swift
//  uzucoin-client
//
//  Created by KAGE on 2018/12/22.
//  Copyright © 2018 kagemiku. All rights reserved.
//

import UIKit

class ViewController: UIViewController {

    override func viewDidLoad() {
        super.viewDidLoad()
        // Do any additional setup after loading the view, typically from a nib.

        let client = Uzucoin_UzucoinServiceClient.init(address: "127.0.0.1:50051", secure: false)
        var request = Uzucoin_RegisterProducerRequest()
        request.uid = "kagemiku"
        request.name = "kage"
        let response = try? client.registerProducer(request)
        print("Response: \(response!.succeeded)")
    }


}

