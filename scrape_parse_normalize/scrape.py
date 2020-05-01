#!/usr/bin/python

import requests
import argparse


HOST = "localhost"
PORT = 8000


def main(args):
    url = 'http://{}:{}/sources/{}/scrape'.format(
        HOST,
        PORT,
        args.source
    )
    params = {}
    if args.max_categories:
        params['max_categories'] = args.max_categories
    if args.max_products_per_category:
        params['max_products_per_category'] = args.max_products_per_category

    result = requests.post(url, params=params).json()
    print('TASK ID: {}'.format(result['task'].get('id')))


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Scrape, Parse, Normalize')
    parser.add_argument("source", help="Source to scrape")
    parser.add_argument(
        "-c", "--max-categories", type=int,
        help="Maximum number of categories to scrape"
    )
    parser.add_argument(
        "-p", "--max-products-per-category", type=int,
        help='Maximum number of products to scrape per category'
    )
    args = parser.parse_args()
    main(args)
