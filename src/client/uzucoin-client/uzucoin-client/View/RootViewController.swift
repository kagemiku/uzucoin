//
//  RootViewController.swift
//  uzucoin-client
//
//  Created by Akira Fukunaga on 2018/12/23.
//  Copyright © 2018 kagemiku. All rights reserved.
//

import UIKit

class RootViewController: UIViewController {

    override func viewDidLoad() {
        super.viewDidLoad()

        let vc = UzucoinTabBarViewController()
        self.addChild(vc)
        self.view.addSubview(vc.view)
        vc.didMove(toParent: self)
    }

    override func viewDidAppear(_ animated: Bool) {
        super.viewDidAppear(animated)

        let userDefaults = UserDefaults.standard
        if let registered = userDefaults.object(forKey: DefaultsKeys.registered.rawValue) as? Bool, !registered {
            let registrationVC = RegistrationViewController()
            self.present(registrationVC, animated: true)
        }
    }

}
