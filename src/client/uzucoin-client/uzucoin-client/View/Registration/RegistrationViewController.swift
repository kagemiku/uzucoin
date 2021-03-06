//
//  RegistrationViewController.swift
//  uzucoin-client
//
//  Created by Akira Fukunaga on 2018/12/23.
//  Copyright © 2018 kagemiku. All rights reserved.
//

import UIKit

protocol RegistrationViewControllerDelegate: class {
    func dismissRegistration()
}

class RegistrationViewController: UIViewController {

    @IBOutlet weak var producerIDLabel: UILabel! {
        didSet {
            self.producerIDLabel.text = self.producerID
        }
    }
    @IBOutlet weak var nameTextField: UITextField!

    private static let idLength = 10
    private let producerID = RegistrationViewController.generateProducerID(length: RegistrationViewController.idLength)

    override func viewDidLoad() {
        super.viewDidLoad()
    }

    private static func generateProducerID(length: Int) -> String {
        let base = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
        var randomString: String = ""

        for _ in 0 ..< length {
            let randomValue = arc4random_uniform(UInt32(base.count))
            randomString += "\(base[base.index(base.startIndex, offsetBy: Int(randomValue))])"
        }

        return randomString
    }

    private func saveProducerID() {
        let userDefaults = UserDefaults.standard
        userDefaults.set(self.producerID, forKey: DefaultsKeys.producerID.rawValue)
        userDefaults.synchronize()
    }
}

extension RegistrationViewController {

    @IBAction func didTapRegister(_ sender: Any) {
        self.nameTextField.resignFirstResponder()
        guard let name = self.nameTextField.text, !name.isEmpty else {
            return
        }

        print("ProducerID: \(self.producerID)")

        var request = Uzucoin_RegisterProducerRequest()
        request.uid = self.producerID
        request.name = name
        let response = try? ProtobufClient.shared.client.registerProducer(request)

        if let res = response, res.succeeded {
            let userDefaults = UserDefaults.standard
            userDefaults.set(true, forKey: DefaultsKeys.registered.rawValue)
            userDefaults.synchronize()
            self.saveProducerID()

            let vc = RegistrationFinishViewController()
            vc.registrationDelegate = self
            vc.modalPresentationStyle = .overCurrentContext
            self.present(vc, animated: true)
        } else {
            print("error")
        }
    }

}

extension RegistrationViewController: RegistrationViewControllerDelegate {

    func dismissRegistration() {
        self.dismiss(animated: true)
    }

}
