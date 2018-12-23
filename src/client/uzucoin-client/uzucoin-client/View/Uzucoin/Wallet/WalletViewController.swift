//
//  WalletViewController.swift
//  uzucoin-client
//
//  Created by Akira Fukunaga on 2018/12/23.
//  Copyright © 2018 kagemiku. All rights reserved.
//

import UIKit

class WalletViewController: UIViewController {

    @IBOutlet weak var balanceLabel: UILabel! {
        didSet {
            self.balanceLabel.text = "0"
        }
    }

    override func viewDidLoad() {
        super.viewDidLoad()

        // Do any additional setup after loading the view.
    }

    override func viewWillAppear(_ animated: Bool) {
        super.viewWillAppear(animated)

        self.updateBalance()
    }

    func updateBalance() {
        guard let producerID = UserDefaults.standard.object(forKey: DefaultsKeys.producerID.rawValue) as? String else { return }
        var request = Uzucoin_GetBalanceRequest()
        request.uid = producerID
        let response = try? ProtobufClient.shared.client.getBalance(request)

        if let res = response {
            let balanceString = String(res.balance)
            DispatchQueue.main.async { [weak self] in
                self?.balanceLabel.text = balanceString
            }
        }
    }

}

extension WalletViewController {

    @IBAction func didTapSendingButton(_ sender: Any) {
        let vc = SendingTopViewController()
        self.navigationController?.pushViewController(vc, animated: true)
    }

}
