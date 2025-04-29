-- phpMyAdmin SQL Dump
-- version 4.7.9
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1
-- Generation Time: Apr 29, 2025 at 03:30 PM
-- Server version: 10.1.31-MariaDB
-- PHP Version: 7.2.3

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `dibimbing_atm`
--

-- --------------------------------------------------------

--
-- Table structure for table `accounts`
--

CREATE TABLE `accounts` (
  `id` int(5) NOT NULL,
  `name` varchar(30) NOT NULL,
  `pin` varchar(6) NOT NULL,
  `balance` decimal(10,0) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `accounts`
--

INSERT INTO `accounts` (`id`, `name`, `pin`, `balance`, `created_at`) VALUES
(1, 'Rifqi', '12345', '0', '2025-04-27 01:47:36'),
(2, 'Rifqi', '54321', '100000', '2025-04-29 10:10:02'),
(3, 'Robi', '12345', '100000', '2025-04-29 10:10:02');

-- --------------------------------------------------------

--
-- Table structure for table `transactions`
--

CREATE TABLE `transactions` (
  `id` int(11) NOT NULL,
  `account_id` int(5) NOT NULL,
  `type` enum('deposit','withdraw','transfer_in','transfer_out') NOT NULL,
  `amount` decimal(10,0) NOT NULL,
  `target_id` int(11) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `transactions`
--

INSERT INTO `transactions` (`id`, `account_id`, `type`, `amount`, `target_id`, `created_at`) VALUES
(1, 2, 'deposit', '100000', 0, '2025-04-27 01:49:48'),
(2, 2, 'withdraw', '50000', 0, '2025-04-29 09:25:12'),
(3, 2, 'deposit', '150000', 0, '2025-04-29 09:25:24'),
(4, 2, 'transfer_out', '100000', 0, '2025-04-29 10:06:03'),
(5, 3, 'transfer_in', '100000', 0, '2025-04-29 10:06:03'),
(6, 3, 'transfer_out', '100000', 0, '2025-04-29 10:08:31'),
(7, 2, 'transfer_in', '100000', 0, '2025-04-29 10:08:31'),
(8, 2, 'transfer_out', '100000', 0, '2025-04-29 10:10:02'),
(9, 3, 'transfer_in', '100000', 0, '2025-04-29 10:10:02');

--
-- Indexes for dumped tables
--

--
-- Indexes for table `accounts`
--
ALTER TABLE `accounts`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `transactions`
--
ALTER TABLE `transactions`
  ADD PRIMARY KEY (`id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `accounts`
--
ALTER TABLE `accounts`
  MODIFY `id` int(5) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- AUTO_INCREMENT for table `transactions`
--
ALTER TABLE `transactions`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=10;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
