package certificate

import (
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/json"
	"errors"
)

var (
	OidOrganization       = []int{2, 5, 4, 10} // 证书类型 (string)
	OidOrganizationalUnit = []int{2, 5, 4, 11} // 标识ID (string)
	OidLocality           = []int{2, 5, 4, 7}  // 扩展ID (string)

	OidVersion     = []int{2, 8, 8, 6, 3, 1, 1} // 内部版本号 (int)
	OidIsTemporary = []int{2, 8, 8, 6, 3, 1, 2} // 是否临时证书 (bool)
	OidBindIP      = []int{2, 8, 8, 6, 3, 1, 3} // 绑定IP ([]string)
	OidNotBefore   = []int{2, 8, 8, 6, 3, 1, 4} // 起始有效期 (time)
	OidNotAfter    = []int{2, 8, 8, 6, 3, 1, 5} // 截止有效期 (time)

	OidUserName     = []int{2, 8, 8, 6, 3, 2, 1} // 用户姓名 (string)
	OidIsSuperAdmin = []int{2, 8, 8, 6, 3, 2, 2} // 是否超级管理员 (bool)

	OidGatewayName  = []int{2, 8, 8, 6, 3, 3, 1} // 网关名称 (string)
	OidMaichineName = []int{2, 8, 8, 6, 3, 3, 2} // 主机名称 (string)
)

type Oid struct {
	Extensions []pkix.Extension
}

func (s *Oid) GetValue(oid asn1.ObjectIdentifier, v interface{}) error {
	exts := s.Extensions
	extLen := len(exts)
	for idx := 0; idx < extLen; idx++ {
		if exts[idx].Id.Equal(oid) {
			err := json.Unmarshal(exts[idx].Value, v)
			if err != nil {
				return err
			} else {
				return nil
			}
		}
	}

	return errors.New("not found")
}
