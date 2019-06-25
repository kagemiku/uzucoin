//
//  SminingHelper.swift
//  uzucoin-client
//
//  Created by Akira Fukunaga on 2019/06/24.
//  Copyright © 2019 kagemiku. All rights reserved.
//

import Foundation
import Metal

class SminingHelper {

    typealias CompletionHandler = (String) -> Void

    static let shared = SminingHelper()
    static private let maxNonceLength = 100

    private let device: MTLDevice
    private let commandQueue: MTLCommandQueue
    private let computePipelineState: MTLComputePipelineState

    private init() {
        device = MTLCreateSystemDefaultDevice()!
        commandQueue = device.makeCommandQueue()!

        let defaultLibrary = device.makeDefaultLibrary()!
        let fun = defaultLibrary.makeFunction(name: "ResolveNonce")!
        computePipelineState = try! device.makeComputePipelineState(function: fun)
    }

    func resolveNonce(with timestamp: String, prevHash: String, completionHander: @escaping CompletionHandler) {
        // command buffer and encoder setting
        let commandBuffer = commandQueue.makeCommandBuffer()!
        let computeCommandEncoder = commandBuffer.makeComputeCommandEncoder()!
        computeCommandEncoder.setComputePipelineState(computePipelineState)

        // input data
        let timestampCString = (timestamp as NSString).utf8String!
        let prevHashCString = (prevHash as NSString).utf8String!

        // input buffers
        let timestampBuffer = device.makeBuffer(bytes: timestampCString, length: timestamp.utf8.count, options: [])!
        let prevHashBuffer = device.makeBuffer(bytes: prevHashCString, length: prevHash.utf8.count, options: [])!
        computeCommandEncoder.setBuffer(timestampBuffer, offset: 0, index: 0)
        computeCommandEncoder.setBuffer(prevHashBuffer, offset: 0, index: 1)

        // output buffer
        let noncePtr = UnsafeMutablePointer<Int8>.allocate(capacity: SminingHelper.maxNonceLength)
        let nonceBuffer = device.makeBuffer(bytes: noncePtr, length: SminingHelper.maxNonceLength, options: [])!
        computeCommandEncoder.setBuffer(nonceBuffer, offset: 0, index: 2)

        // number of threads per group and thread groups
        let width = 64
        let threadsPerGroup = MTLSize(width: /*width*/ 1, height: 1, depth: 1)
        let threadgroupsPerGrid = MTLSize(width: /*(timestamp.utf8.count + width - 1) / width*/ 1, height: 1, depth: 1)
        computeCommandEncoder.dispatchThreadgroups(threadgroupsPerGrid, threadsPerThreadgroup: threadsPerGroup)

        computeCommandEncoder.endEncoding()

        commandBuffer.addCompletedHandler { _ in
            let rawData = Data(bytesNoCopy: nonceBuffer.contents(), count: SminingHelper.maxNonceLength, deallocator: .none)
            let resultPtr: UnsafePointer<Int8> = rawData.withUnsafeBytes {
                $0
            }

            let nonceString = String(cString: resultPtr)
            completionHander(nonceString)
        }

        commandBuffer.commit()
    }
}
