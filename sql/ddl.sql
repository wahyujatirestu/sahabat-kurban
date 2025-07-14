-- Enable pgcrypto extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Fungsi trigger untuk update updated_at otomatis
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Enum status penerima daging
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_penerima_enum') THEN
        CREATE TYPE status_penerima_enum AS ENUM ('warga', 'dhuafa', 'panitia', 'pekurban');
    END IF;
END$$;

-- Tabel users (admin, panitia, user biasa)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'panitia', 'user')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

CREATE TRIGGER trigger_update_users
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Tabel pekurban (boleh terkait user atau tidak)
CREATE TABLE pekurban (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE NULL,
    name VARCHAR(100) NULL,
    phone VARCHAR(20),
    email VARCHAR(100),
    alamat TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE TRIGGER trigger_update_pekurban
BEFORE UPDATE ON pekurban
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Enum jenis hewan
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'jenis_hewan_enum') THEN
        CREATE TYPE jenis_hewan_enum AS ENUM ('sapi', 'kambing', 'domba');
    END IF;
END$$;

-- Tabel hewan_kurban
CREATE TABLE hewan_kurban (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    jenis jenis_hewan_enum NOT NULL,
    berat NUMERIC(5,2) NOT NULL,
    tanggal_pendaftaran DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

CREATE TRIGGER trigger_update_hewan_kurban
BEFORE UPDATE ON hewan_kurban
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Tabel pekurban_hewan 
CREATE TABLE pekurban_hewan (
    pekurban_id UUID NOT NULL,
    hewan_id UUID NOT NULL,
    porsi NUMERIC(4,3) NOT NULL CHECK (porsi > 0 AND porsi <= 1),
    PRIMARY KEY (pekurban_id, hewan_id),
    FOREIGN KEY (pekurban_id) REFERENCES pekurban(id) ON DELETE CASCADE,
    FOREIGN KEY (hewan_id) REFERENCES hewan_kurban(id) ON DELETE CASCADE
);

-- Tabel penyembelihan
CREATE TABLE penyembelihan (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hewan_id UUID NOT NULL UNIQUE,
    tanggal_penyembelihan DATE NOT NULL,
    lokasi VARCHAR(255),
    urutan_rencana INT DEFAULT 9999 NOT NULL,
    urutan_aktual INT DEFAULT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    FOREIGN KEY (hewan_id) REFERENCES hewan_kurban(id) ON DELETE CASCADE
);

CREATE TRIGGER trigger_update_penyembelihan
BEFORE UPDATE ON penyembelihan
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Tabel penerima_daging dengan relasi ke pekurban
CREATE TABLE penerima_daging (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    alamat TEXT,
    phone VARCHAR(20),
    status status_penerima_enum NOT NULL DEFAULT 'warga',
    pekurban_id UUID UNIQUE NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    FOREIGN KEY (pekurban_id) REFERENCES pekurban(id) ON DELETE SET NULL
);

CREATE TRIGGER trigger_update_penerima_daging
BEFORE UPDATE ON penerima_daging
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Tabel distribusi_daging
CREATE TABLE distribusi_daging (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    penerima_id UUID UNIQUE NOT NULL,
    hewan_id UUID NOT NULL,
    jumlah_paket INT NOT NULL CHECK (jumlah_paket > 0),
    tanggal_distribusi DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    FOREIGN KEY (penerima_id) REFERENCES penerima_daging(id) ON DELETE CASCADE,
    FOREIGN KEY (hewan_id) REFERENCES hewan_kurban(id) ON DELETE CASCADE
);

CREATE TRIGGER trigger_update_distribusi_daging
BEFORE UPDATE ON distribusi_daging
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Tabel pembayaran_kurban (fleksibel metode dan status)
CREATE TABLE pembayaran_kurban (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pekurban_id UUID NOT NULL,
    metode VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    jumlah NUMERIC(12,2) NOT NULL,
    tanggal_pembayaran TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    payment_url TEXT,
    transaction_id VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    FOREIGN KEY (pekurban_id) REFERENCES pekurban(id) ON DELETE CASCADE
);

CREATE TRIGGER trigger_update_pembayaran_kurban
BEFORE UPDATE ON pembayaran_kurban
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Tabel refresh_tokens untuk refresh token JWT
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked BOOLEAN DEFAULT FALSE NOT NULL,
    revoked_at TIMESTAMP WITH TIME ZONE,
    replaced_by_token TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TRIGGER trigger_update_refresh_tokens
BEFORE UPDATE ON refresh_tokens
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
