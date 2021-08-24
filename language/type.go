package language

import (
	"reflect"
	"strings"

	"xorm.io/xorm/schemas"
)

func SQLType2GoType(st schemas.SQLType) reflect.Type {
	name := strings.ToUpper(st.Name)
	switch name {
	case schemas.Bit, schemas.SmallInt, schemas.MediumInt, schemas.Int, schemas.Integer, schemas.Serial:
		return schemas.IntType
	case schemas.UnsignedBit, schemas.UnsignedTinyInt, schemas.UnsignedInt, schemas.UnsignedSmallInt, schemas.UnsignedMediumInt:
		return reflect.TypeOf((*uint)(nil)).Elem()
	case schemas.TinyInt:
		if st.DefaultLength == 1 {
			return schemas.BoolType
		}
		return schemas.IntType
	case schemas.BigInt, schemas.BigSerial:
		return schemas.Int64Type
	case schemas.UnsignedBigInt:
		return reflect.TypeOf((*uint64)(nil)).Elem()
	case schemas.Float, schemas.Real:
		return schemas.Float32Type
	case schemas.Double:
		return schemas.Float64Type
	case schemas.Char, schemas.NChar, schemas.Varchar, schemas.NVarchar, schemas.TinyText, schemas.Text, schemas.NText, schemas.MediumText, schemas.LongText, schemas.Enum, schemas.Set, schemas.Uuid, schemas.Clob, schemas.SysName:
		return schemas.StringType
	case schemas.TinyBlob, schemas.Blob, schemas.LongBlob, schemas.Bytea, schemas.Binary, schemas.MediumBlob, schemas.VarBinary, schemas.UniqueIdentifier:
		return schemas.BytesType
	case schemas.Bool:
		return schemas.BoolType
	case schemas.DateTime, schemas.Date, schemas.Time, schemas.TimeStamp, schemas.TimeStampz, schemas.SmallDateTime, schemas.Year:
		return schemas.TimeType
	case schemas.Decimal, schemas.Numeric, schemas.Money, schemas.SmallMoney:
		return schemas.StringType
	default:
		return schemas.StringType
	}
}
