#!/usr/bin/env python
# -*- coding: utf-8 -*-
import sys
from urllib import urlencode, urlopen

__version__ = '0.2.0'
help = """
%s

program [status | login | logout]

status          show login status
login           login
logout          logout
""" % __version__


class MD6(object):
    def mc(self, a):
        ret = ''
        b = '0123456789ABCDEF'
        if a == ord(' '):
            ret = '+'
        elif ((a < ord('0') and a != ord('-')
               and a != ord('.'))
              or (a < ord('A') and a > ord('9'))
              or (a > ord('Z') and a < ord('a')
                  and a != ord('_'))
              or (a > ord('z'))
              ):
            ret = '%'
            ret += b[a >> 4]
            ret += b[a & 15]
        else:
            ret = chr(a)
        return ret

    def m(self, a):
        return (((a & 1) << 7)
                | ((a & (0x2)) << 5)
                | ((a & (0x4)) << 3)
                | ((a & (0x8)) << 1)
                | ((a & (0x10)) >> 1)
                | ((a & (0x20)) >> 3)
                | ((a & (0x40)) >> 5)
                | ((a & (0x80)) >> 7)
                )

    def md6(self, s):
        b = ''
        c = 0xbb
        for i in range(len(s)):
            c = self.m(ord(s[i])) ^ (0x35 ^ (i & 0xff))
            b += self.mc(c)
        return b

    def __call__(self, s):
        return self.md6(s)


class Account(object):
    url_test = 'http://www.baidu.com'
    url_login = 'http://8.8.8.8:90/login'
    start = 'http://8.8.8.8:90'

    def status(self, *args, **kwargs):
        r = urlopen(self.url_test)
        redirect = r.url.startswith(self.start)
        if redirect:
            print('logged out')
        else:
            print('logged in')

    def login(self, username, pasword):
        # 设置表单参数
        url = urlopen(self.url_login).url
        if url == self.url_login:
            print('already logged in')
            return

        uri = url.split('?')[1]
        data = urlencode(dict(
            username=username,
            password=password,
            uri=uri,
            terminal='pc',
            login_type='login',
            check_passwd=0,
            show_tip='block',
            show_change_password='block',
            short_message='none',
            show_captcha='none',
            show_read='block',
            show_assure='none',
            assure_phone='',
            password1='',
            new_password='',
            retype_newpassword='',
            captcha_value='',
            save_user=1,
            save_pass=1,
            read=1
        ))

        r = urlopen(self.url_login, data=data)
        text = r.read()
        if username in text:
            print('login success')
        else:
            print('fail')
            print('username: %s' % username)
            print('password: %s' % password)
            print('data: %s' % data)
            print(text)

    def logout(self, *args, **kwargs):
        data = urlencode({
            'login_type': 'logout',
        })
        urlopen(self.url_login, data=data)
        self.status()

if __name__ == '__main__':
    args = sys.argv[1:]
    account = Account()

    funcs = {
        'status': account.status,
        'login': account.login,
        'logout': account.logout
    }
    username = "username"
    password = "password"
    password = MD6()(password)

    if not args or len(args) > 1:
        print(help)
        sys.exit()

    arg = args[0]
    if arg not in funcs:
        print(help)
    else:
        funcs[arg](username, password)
