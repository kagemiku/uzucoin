//
//  UzucoinTabBarViewController.swift
//  uzucoin-client
//
//  Created by Akira Fukunaga on 2018/12/23.
//  Copyright © 2018 kagemiku. All rights reserved.
//

import UIKit

class UzucoinTabBarViewController: UITabBarController {

    private lazy var walletVC: UINavigationController = {
        let vc = WalletViewController()
        let nvc = UINavigationController(rootViewController: vc)
        nvc.title = "Wallet"
        return nvc
    }()

    private lazy var historyVC: UINavigationController = {
        let vc = HistoryViewController()
        let nvc = UINavigationController(rootViewController: vc)
        nvc.title = "履歴"
        return nvc
    }()

    private lazy var sminingVC: UINavigationController = {
        let vc = UIViewController()
        let nvc = UINavigationController(rootViewController: vc)
        nvc.title = "S(min)ing!"
        return nvc
    }()

    override func viewDidLoad() {
        super.viewDidLoad()

        self.setViewControllers(
            [
                self.walletVC,
                self.historyVC,
                self.sminingVC,
            ],
            animated: true
        )
    }
}
