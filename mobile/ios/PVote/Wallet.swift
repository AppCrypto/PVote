import Foundation
import Security

enum Wallet {
    private static let service = "org.appcrypto.pvote.research.wallet"
    private static let account = "voter"
    private static var query: [String: Any] {
        [kSecClass as String: kSecClassGenericPassword, kSecAttrService as String: service,
         kSecAttrAccount as String: account]
    }
    static func load() throws -> String? {
        var q = query
        q[kSecReturnData as String] = true
        var result: CFTypeRef?
        let status = SecItemCopyMatching(q as CFDictionary, &result)
        if status == errSecItemNotFound { return nil }
        guard status == errSecSuccess, let data = result as? Data,
              let key = String(data: data, encoding: .utf8) else { throw PVoteError("Cannot read local wallet") }
        return key
    }
    static func save(_ key: String) throws {
        let data = Data(key.utf8)
        let status = SecItemUpdate(query as CFDictionary, [kSecValueData as String: data] as CFDictionary)
        if status == errSecSuccess { return }
        guard status == errSecItemNotFound else { throw PVoteError("Cannot update local wallet") }
        var q = query
        q[kSecValueData as String] = data
        q[kSecAttrAccessible as String] = kSecAttrAccessibleWhenUnlockedThisDeviceOnly
        guard SecItemAdd(q as CFDictionary, nil) == errSecSuccess else { throw PVoteError("Cannot save local wallet") }
    }
}
