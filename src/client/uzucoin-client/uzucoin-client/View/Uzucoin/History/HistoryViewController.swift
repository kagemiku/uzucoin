//
//  HistoryViewController.swift
//  uzucoin-client
//
//  Created by Akira Fukunaga on 2018/12/23.
//  Copyright © 2018 kagemiku. All rights reserved.
//

import UIKit

class HistoryViewController: UIViewController {

    @IBOutlet weak var historyTableView: UITableView! {
        didSet {
            self.historyTableView.dataSource = self
            self.historyTableView.delegate = self
            self.historyTableView.register(UINib(nibName: "HistoryCellView", bundle: nil), forCellReuseIdentifier: HistoryCellView.reuseIdentifier)
        }
    }

    private var data: [Uzucoin_Transaction] = [] {
        didSet {
            self.historyTableView.reloadData()
        }
    }

    override func viewDidLoad() {
        super.viewDidLoad()
    }

    override func viewWillAppear(_ animated: Bool) {
        super.viewWillAppear(animated)

        self.getHistory()
    }

    private func getHistory() {
        guard let uid = UserDefaults.standard.object(forKey: DefaultsKeys.producerID.rawValue) as? String else { return }
        var request = Uzucoin_GetHistoryRequest()
        request.uid = uid
        let response: Uzucoin_History
        do {
            response = try ProtobufClient.shared.client.getHistory(request)
        } catch (let error) {
            print(error)
            return
        }

        self.data = response.transactions.reversed()
    }

}

extension HistoryViewController: UITableViewDataSource {

    func tableView(_ tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
        return self.data.count
    }

    func tableView(_ tableView: UITableView, cellForRowAt indexPath: IndexPath) -> UITableViewCell {
        let cell = tableView.dequeueReusableCell(withIdentifier: HistoryCellView.reuseIdentifier, for: indexPath) as! HistoryCellView
        let cellData = self.data[indexPath.row]
        cell.configure(dateString: cellData.timestamp, name: cellData.toUid, amount: cellData.amount)

        return cell
    }

}

extension HistoryViewController: UITableViewDelegate {

}
