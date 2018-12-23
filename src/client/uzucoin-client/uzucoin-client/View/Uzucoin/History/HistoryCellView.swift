//
//  HistoryCellView.swift
//  uzucoin-client
//
//  Created by Akira Fukunaga on 2018/12/23.
//  Copyright © 2018 kagemiku. All rights reserved.
//

import UIKit

class HistoryCellView: UITableViewCell {

    @IBOutlet weak var dateLabel: UILabel!
    @IBOutlet weak var nameLabel: UILabel!
    @IBOutlet weak var amountLabel: UILabel!

    override func awakeFromNib() {
        super.awakeFromNib()
        // Initialization code
    }

    override func setSelected(_ selected: Bool, animated: Bool) {
        super.setSelected(selected, animated: animated)

        // Configure the view for the selected state
    }

    override func prepareForReuse() {
        super.prepareForReuse()

        self.dateLabel.text = ""
        self.nameLabel.text = ""
        self.amountLabel.text = ""
    }

    func configure(dateString: String, name: String, amount: Double) {
        self.dateLabel.text = dateString
        self.nameLabel.text = name
        self.amountLabel.text = String(amount)
    }

}
