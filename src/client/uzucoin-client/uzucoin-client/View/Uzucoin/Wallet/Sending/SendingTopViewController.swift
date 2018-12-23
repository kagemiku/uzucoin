//
//  SendingTopViewController.swift
//  uzucoin-client
//
//  Created by Akira Fukunaga on 2018/12/23.
//  Copyright © 2018 kagemiku. All rights reserved.
//

import UIKit

class SendingTopViewController: UIViewController {

    @IBOutlet weak var destProducerIDTextField: UITextField!
    @IBOutlet weak var amountTextField: UITextField!

    override func viewDidLoad() {
        super.viewDidLoad()

        self.navigationItem.title = "送金"
    }

}

extension SendingTopViewController {

    @IBAction func didTapSendButton(_ sender: Any) {
        self.destProducerIDTextField.resignFirstResponder()
        self.amountTextField.resignFirstResponder()

        guard
            let fromProducerID = UserDefaults.standard.object(forKey: DefaultsKeys.producerID.rawValue) as? String,
            let destProducerID = self.destProducerIDTextField.text?.trimmingCharacters(in: CharacterSet.whitespaces),
            let amountText = self.amountTextField.text,
            let amount = Double(amountText)
        else {
            return
        }

        var request = Uzucoin_AddTransactionRequest()
        request.fromUid = fromProducerID
        request.toUid = destProducerID
        request.amount = amount
        let response: Uzucoin_AddTransactionResponse
        do {
            response = try ProtobufClient.shared.client.addTransaction(request)
        } catch (let error) {
            print(error)
            return
        }

        if !response.timestamp.isEmpty {
            self.navigationController?.popToRootViewController(animated: true)
        }
    }
}
