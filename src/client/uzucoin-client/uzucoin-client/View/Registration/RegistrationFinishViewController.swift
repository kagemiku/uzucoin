//
//  RegistrationFinishViewController.swift
//  uzucoin-client
//
//  Created by Akira Fukunaga on 2018/12/23.
//  Copyright © 2018 kagemiku. All rights reserved.
//

import UIKit

class RegistrationFinishViewController: UIViewController {

    var registrationDelegate: RegistrationViewControllerDelegate? = nil

    override func viewDidLoad() {
        super.viewDidLoad()

        // Do any additional setup after loading the view.
    }

}

extension RegistrationFinishViewController {

    @IBAction func didTapHomeButton(_ sender: Any) {
        self.dismiss(animated: true) { [weak self] in
            guard let delegate = self?.registrationDelegate else { return }
            delegate.dismissRegistration()
        }
    }

}
