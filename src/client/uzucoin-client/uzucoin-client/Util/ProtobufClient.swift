//
//  ProtobufClient.swift
//  uzucoin-client
//
//  Created by Akira Fukunaga on 2018/12/23.
//  Copyright © 2018 kagemiku. All rights reserved.
//

import Foundation

final class ProtobufClient {
    static let shared = ProtobufClient.init()

    private static let host = "127.0.0.1"
    private static let port = "50051"
    private static let address = "\(ProtobufClient.host):\(ProtobufClient.port)"
    let client = Uzucoin_UzucoinServiceClient.init(address: ProtobufClient.address, secure: false)

    private init() { }
}
