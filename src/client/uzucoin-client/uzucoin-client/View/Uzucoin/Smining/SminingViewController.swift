//
//  SminingViewController.swift
//  uzucoin-client
//
//  Created by Akira Fukunaga on 2019/12/07.
//  Copyright © 2019 kagemiku. All rights reserved.
//

import UIKit

class SminingViewController: UIViewController {

    private var latestPrevHash = ""

    // MARK: IBOutlets

    @IBOutlet private weak var taskLabel: UILabel!
    @IBOutlet private weak var nonceTextField: UITextField!
    @IBOutlet private weak var sendButton: UIButton!

    override func viewDidLoad() {
        super.viewDidLoad()

        resetLabels()
    }

    override func viewWillAppear(_ animated: Bool) {
        super.viewWillAppear(animated)

        updateTask()
    }

    private func resetLabels() {
        nonceTextField.text = ""
    }

    private func updateTask() {
        _ = try? ProtobufClient.shared.client.getTask(.init()) { [weak self] (task, result) in
            guard let task = task else { print("error"); return }

            let taskLabelText: String
            if task.exists {
                taskLabelText = task.transaction.timestamp
                self?.latestPrevHash = task.prevHash
            } else {
                taskLabelText = "None"
            }

            DispatchQueue.main.async {
                self?.taskLabel.text = taskLabelText
                self?.nonceTextField.isEnabled = task.exists
                self?.sendButton.isEnabled = task.exists
            }
        }
    }

    private func sendNonce(_ nonce: String) {
        guard let producerID = UserDefaults.standard.object(forKey: DefaultsKeys.producerID.rawValue) as? String else { return }

        var request = Uzucoin_ResolveNonceRequest()
        request.nonce = nonce
        request.resolverUid = producerID
        request.prevHash = latestPrevHash
        _ = try? ProtobufClient.shared.client.resolveNonce(request) { [weak self] (response, result) in
            guard let response = response else { print("error"); return }

            print(response.reward)
            print(response.succeeded)
            self?.updateTask()

            DispatchQueue.main.async {
                self?.showAlert(response.succeeded)
                self?.resetLabels()
            }
        }
    }

    private func showAlert(_ success: Bool) {
        let alert: UIAlertController
        if success {
            alert = UIAlertController(title: "S(min)ING!", message: "成功です♪ 🎉🎉🎉🎉🎉🎉🎉🎉", preferredStyle: .alert)
            alert.addAction(.init(title: "OK", style: .default))
        } else {
            alert = UIAlertController(title: "S(min)ING!", message: "失敗です。。。", preferredStyle: .alert)
            alert.addAction(.init(title: "OK", style: .default))
        }

        navigationController?.present(alert, animated: true)
    }
}

extension SminingViewController {

    @IBAction func didTapSendButton(_ sender: Any) {
        guard let nonce = nonceTextField.text else { return }

        sendNonce(nonce)
    }

}
