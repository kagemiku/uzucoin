//
//  SminingViewController.swift
//  uzucoin-client
//
//  Created by Akira Fukunaga on 2019/06/24.
//  Copyright © 2019 kagemiku. All rights reserved.
//

import UIKit

class SminingViewController: UIViewController {

    override func viewDidLoad() {
        super.viewDidLoad()
    }

    @IBAction func didTapSminingButton(_ sender: Any) {
        let timestamp = "2019/06/24"
        let prevHash = "Uzuki"

        SminingHelper.shared.resolveNonce(with: timestamp, prevHash: prevHash) { (nonce) in
            print(nonce)
        }
    }
    
}
